package vk

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/gobwas/vk/internal/httputil"
	"github.com/gobwas/vk/internal/syncutil"
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

type Iterator struct {
	Method  string
	Options []QueryOption
	Parse   func([]byte) (int, error)

	limiter *syncutil.Limiter
	offset  int
	err     error
	once    sync.Once
}

func (it *Iterator) Next(ctx context.Context) bool {
	if it.err != nil {
		return false
	}

	it.init()

	var (
		n   int
		err error
	)
	for {
		it.limiter.Do(func() {
			n, err = it.fetch(ctx)
		})
		if vkErr, ok := err.(*Error); ok && vkErr.Temporary() {
			continue
		}
		break
	}

	it.err = err

	return err == nil && n > 0
}

func (it *Iterator) Err() error {
	return it.err
}

func (it *Iterator) init() {
	it.once.Do(func() {
		it.limiter = syncutil.NewLimiter(time.Second, 3)
	})
}

func (it *Iterator) Close() {
	it.init()
	it.limiter.Close()
}

func (it *Iterator) fetch(ctx context.Context) (int, error) {
	bts, err := Request(ctx, it.Method,
		WithOptions(it.Options),
		WithNumber("offset", it.offset),
	)
	if err != nil {
		return 0, err
	}
	bts, err = StripResponse(bts)
	if err != nil {
		return 0, err
	}
	n, err := it.Parse(bts)
	if err != nil {
		return 0, err
	}

	it.offset += n

	return n, nil
}
