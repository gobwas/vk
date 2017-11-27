package vk

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
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
	State        string
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

func (a *App) AuthorizePath(redirect string) string {
	auth, err := url.Parse("https://oauth.vk.com/authorize")
	if err != nil {
		panic("constant url is invalid: " + err.Error())
	}
	query := url.Values{
		"v":             []string{version},
		"client_id":     []string{a.ClientID},
		"state":         []string{a.State},
		"redirect_uri":  []string{redirect},
		"scope":         []string{a.Scope.String()},
		"display":       []string{"page"},
		"response_type": []string{"code"},
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

	expires := time.Now().Add(time.Second * time.Duration(acc.Expires))

	return &AccessToken{
		Token:   acc.Token,
		UserID:  acc.UserID,
		Expires: expires,
	}, nil
}

func CodeFromParams(params url.Values) (string, error) {
	if errName := params.Get("error"); errName != "" {
		return "", fmt.Errorf(
			"bad redirect: %s: %s",
			errName, params.Get("error_description"),
		)
	}
	return params.Get("code"), nil
}
