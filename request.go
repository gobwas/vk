package vk

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/gobwas/vk/internal/httputil"
)

var (
	DefaultRateInterval = time.Second / 3
	DefaultRateBurst    = 3
)

type Authorizer interface {
	Authorize() (token string, err error)
}

type request struct {
	method string
	query  url.Values
}

type QueryOption func(url.Values)

func QueryOptions(options ...QueryOption) []QueryOption {
	return options
}

func WithOptions(options []QueryOption) QueryOption {
	return func(query url.Values) {
		for _, option := range options {
			option(query)
		}
	}
}

func WithQuery(source url.Values) QueryOption {
	return func(query url.Values) {
		for key, values := range source {
			for _, value := range values {
				query.Add(key, value)
			}
		}
	}
}

func WithParam(key, value string) QueryOption {
	return func(query url.Values) {
		query.Add(key, value)
	}
}

func WithNumber(key string, n int) QueryOption {
	return func(query url.Values) {
		query.Set(key, strconv.Itoa(n))
	}
}

func WithNumbers(key string, ns ...int) QueryOption {
	strs := make([]string, len(ns))
	for i, n := range ns {
		strs[i] = strconv.Itoa(n)
	}
	list := strings.Join(strs, ",")
	return func(query url.Values) {
		query.Set(key, list)
	}
}

func WithStrings(key string, strs ...string) QueryOption {
	list := strings.Join(strs, ",")
	return func(query url.Values) {
		query.Set(key, list)
	}
}

func WithAccessToken(access *AccessToken) QueryOption {
	return func(query url.Values) {
		query.Set("access_token", access.Token)
	}
}

func Request(ctx context.Context, method string, options ...QueryOption) ([]byte, error) {
	req := &request{
		method: method,
		query:  make(url.Values),
	}
	req.query.Set("v", version)

	WithOptions(options)(req.query)

	r, err := http.NewRequest("GET", req.url(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(r.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := httputil.CheckResponseStatus(resp); err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

func StripResponse(p []byte) ([]byte, error) {
	var response Response
	if err := response.UnmarshalJSON(p); err != nil {
		return nil, err
	}
	if err := response.Error(); err != nil {
		return nil, err
	}
	return response.Body, nil
}

func (req *request) url() string {
	u, err := url.Parse("https://api.vk.com/method")
	if err != nil {
		panic("constant url is invalid: " + err.Error())
	}
	u.Path += "/" + req.method
	u.RawQuery = req.query.Encode()
	return u.String()
}

type Caller struct {
	Method         string
	Options        []QueryOption
	Limiter        *rate.Limiter
	ResolveCaptcha func(ctx context.Context, img string) (text string, err error)

	runtime []QueryOption
}

func (c *Caller) Call(ctx context.Context, opts ...QueryOption) ([]byte, error) {
retry:
	if lim := c.Limiter; lim != nil {
		if err := lim.Wait(ctx); err != nil {
			return nil, err
		}
	}
	bts, err := Request(ctx, c.Method,
		WithOptions(c.Options),
		WithOptions(c.runtime),
		WithOptions(opts),
	)
	if err == nil {
		bts, err = StripResponse(bts)
	}
	if c.Limiter != nil && TemporaryError(err) {
		goto retry
	}
	if captcha := c.ResolveCaptcha; captcha != nil {
		sid, img, ok := CaptchaError(err)
		if ok {
			text, rerr := c.ResolveCaptcha(ctx, img)
			if rerr == nil {
				c.runtime = append(c.runtime,
					WithParam("captcha_sid", sid),
					WithParam("captcha_key", text),
				)
				goto retry
			}
		}
	}
	return bts, err
}

type Iterator struct {
	// Caller fields
	Method  string
	Options []QueryOption
	Limiter *rate.Limiter

	Parse func([]byte) (int, error)

	once   sync.Once
	caller Caller
	offset int
	err    error
}

func (it *Iterator) Next(ctx context.Context) bool {
	if it.err != nil {
		return false
	}

	it.init()

	bts, err := it.caller.Call(ctx,
		WithNumber("offset", it.offset),
	)
	if err != nil {
		it.err = err
		return false
	}

	n, err := it.Parse(bts)
	if err != nil {
		it.err = err
		return false
	}
	if n == 0 {
		return false
	}
	it.offset += n

	return true
}

func (it *Iterator) Err() error {
	return it.err
}

func (it *Iterator) init() {
	it.once.Do(func() {
		if it.Limiter == nil {
			it.Limiter = DefaultLimiter()
		}
		it.caller = Caller{
			Method:  it.Method,
			Options: it.Options,
			Limiter: it.Limiter,
		}
	})
}

func DefaultLimiter() *rate.Limiter {
	return rate.NewLimiter(
		rate.Every(DefaultRateInterval),
		DefaultRateBurst,
	)
}

func TemporaryError(err error) bool {
	if vkErr, ok := err.(*Error); ok {
		return vkErr.Temporary()
	}
	return false
}

func CaptchaError(err error) (sid, img string, ok bool) {
	if vkErr, ok := err.(*Error); ok {
		return vkErr.CaptchaSID, vkErr.CaptchaImg, vkErr.CaptchaSID != ""
	}
	return "", "", false
}
