package vk

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gobwas/vk/internal/httputil"
)

type Authorizer interface {
	Authorize() (token string, err error)
}

type request struct {
	method string
	query  url.Values
}

type QueryOption func(url.Values)

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

func WithOffset(offset int) QueryOption {
	return func(query url.Values) {
		query.Set("offset", strconv.Itoa(offset))
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
