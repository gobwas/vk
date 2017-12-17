package vk

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gobwas/vk/internal/httputil"
)

type AccessToken struct {
	Token   string
	Expires time.Time
	UserID  int
}

type App struct {
	ClientID     string
	ClientSecret string
	Scope        Scope
}

func (a *App) Authorize(ctx context.Context, redirectPath, code string) (*AccessToken, error) {
	access, err := http.NewRequest(
		"GET", a.AccessTokenPath(redirectPath, code), nil,
	)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(access.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := httputil.CheckResponseStatus(resp); err != nil {
		return nil, err
	}
	return parseAccessTokenResponse(resp.Body)
}

func (a *App) AuthPathToken(redirect string, options ...QueryOption) string {
	return a.authorizePath(redirect, "token", options...)
}

func (a *App) AuthPathCode(redirect string, options ...QueryOption) string {
	return a.authorizePath(redirect, "code", options...)
}

func (a *App) authorizePath(redirect, authType string, options ...QueryOption) string {
	auth, err := url.Parse("https://oauth.vk.com/authorize")
	if err != nil {
		panic("constant url is invalid: " + err.Error())
	}
	query := url.Values{
		"v":             []string{version},
		"client_id":     []string{a.ClientID},
		"redirect_uri":  []string{redirect},
		"scope":         []string{a.Scope.String()},
		"response_type": []string{authType},
	}
	for _, option := range options {
		option(query)
	}
	auth.RawQuery = query.Encode()
	return auth.String()
}

func (a *App) AccessTokenPath(redirect, code string) string {
	access, err := url.Parse("https://oauth.vk.com/access_token")
	if err != nil {
		panic("constant url is invalid: " + err.Error())
	}
	query := url.Values{
		"client_id":     []string{a.ClientID},
		"client_secret": []string{a.ClientSecret},
		"redirect_uri":  []string{redirect},
		"code":          []string{code},
	}
	access.RawQuery = query.Encode()
	return access.String()
}

func parseAccessTokenResponse(resp io.Reader) (*AccessToken, error) {
	bts, err := ioutil.ReadAll(resp)
	if err != nil {
		return nil, err
	}
	var acc rawAccess
	if err := acc.UnmarshalJSON(bts); err != nil {
		return nil, err
	}

	return &AccessToken{
		Token:   acc.Token,
		UserID:  acc.UserID,
		Expires: expiresDate(acc.Expires),
	}, nil
}

func RedirectQueryError(query url.Values) error {
	errName := query.Get("error")
	if errName == "" {
		return nil
	}
	return fmt.Errorf(
		"bad redirect: %s: %s",
		errName, query.Get("error_description"),
	)
}

func expiresDate(sec int) time.Time {
	return time.Now().Add(time.Second * time.Duration(sec))
}

func CodeFromQuery(query url.Values) (string, error) {
	if err := RedirectQueryError(query); err != nil {
		return "", err
	}
	return query.Get("code"), nil
}

func TokenFromURL(str string) (*AccessToken, error) {
	u, err := url.Parse(str)
	if err != nil {
		return nil, err
	}
	params, err := url.ParseQuery(u.Fragment)
	if err != nil {
		return nil, err
	}
	return TokenFromQuery(params)
}

func TokenFromQuery(query url.Values) (*AccessToken, error) {
	if err := RedirectQueryError(query); err != nil {
		return nil, err
	}
	expires, err := strconv.Atoi(query.Get("expires_in"))
	if err != nil {
		return nil, err
	}
	userID, err := strconv.Atoi(query.Get("user_id"))
	if err != nil {
		return nil, err
	}
	return &AccessToken{
		Token:   query.Get("access_token"),
		UserID:  userID,
		Expires: expiresDate(expires),
	}, nil
}
