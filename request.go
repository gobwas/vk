package vk

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Authorizer interface {
	Authorize() (token string, err error)
}

type request struct {
	query url.Values
}

type RequestOption func(*request)

func WithParam(key, value string) RequestOption {
	return func(r *request) {
		r.query.Add(key, value)
	}
}

func WithAccess(access *Access) RequestOption {
	return func(r *request) {
		r.query.Set("access_token", access.Token)
	}
}

func Request(ctx context.Context, method string, options ...RequestOption) ([]byte, error) {
	req := &request{
		query: make(url.Values),
	}
	req.query.Set("v", version)
	for _, opt := range options {
		opt(req)
	}

	u, err := url.Parse("https://api.vk.com/method")
	if err != nil {
		return nil, err
	}
	u.Path += "/" + method
	u.RawQuery = req.query.Encode()

	r, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(r.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := checkResponseStatus(resp); err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}
