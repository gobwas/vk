package vk

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gobwas/vk/httputil"
)

type Authorizer interface {
	Authorize() (token string, err error)
}

type request struct {
	method string
	query  url.Values
}

type RequestOption func(*request)

func WithParam(key, value string) RequestOption {
	return func(r *request) {
		r.query.Add(key, value)
	}
}

func WithAccessToken(access *AccessToken) RequestOption {
	return func(r *request) {
		r.query.Set("access_token", access.Token)
	}
}

func Request(ctx context.Context, method string, options ...RequestOption) ([]byte, error) {
	req := &request{
		method: method,
		query:  make(url.Values),
	}
	req.query.Set("v", version)
	for _, opt := range options {
		opt(req)
	}
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

func (req *request) url() string {
	u, err := url.Parse("https://api.vk.com/method")
	if err != nil {
		panic("constant url is invalid: " + err.Error())
	}
	u.Path += "/" + req.method
	u.RawQuery = req.query.Encode()
	return u.String()
}
