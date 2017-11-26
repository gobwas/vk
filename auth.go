package vk

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gobwas/vk/internal/httputil"
)

var DefaultRedirectHost = "127.0.0.1"

type AccessToken struct {
	Token   string
	Expires time.Time
	UserID  int
}

type rawAccess struct {
	Token   string `json:"access_token"`
	Expires int    `json:"expires_in"`
	UserID  int    `json:"user_id"`
}

type Auth struct {
	ClientID     string
	ClientSecret string
	State        string
	Scope        Scope

	RedirectHost string
	RedirectPort int
}

func (a *Auth) Authorize(ctx context.Context) (*AccessToken, error) {
	redirect := make(chan requestAndError, 1)
	redirectPath, err := a.redirectServer(redirect)
	if err != nil {
		return nil, err
	}

	// Open a web browser to authorize an app.
	if err := browse(ctx, a.authPath(redirectPath)); err != nil {
		return nil, err
	}
	req, err := waitRedirect(ctx, redirect)
	if err != nil {
		return nil, err
	}
	code, err := codeFromParams(req.URL.Query())
	if err != nil {
		return nil, err
	}

	access, err := http.NewRequest(
		"GET", a.accessTokenPath(redirectPath, code), nil,
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

	return parseAccessTokenResponse(resp)
}

func parseAccessTokenResponse(resp *http.Response) (*AccessToken, error) {
	var acc rawAccess
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&acc); err != nil {
		return nil, err
	}

	expires := time.Now().Add(time.Second * time.Duration(acc.Expires))

	return &AccessToken{
		Token:   acc.Token,
		UserID:  acc.UserID,
		Expires: expires,
	}, nil
}

type requestAndError struct {
	req *http.Request
	err error
}

func (a *Auth) authPath(redirect string) string {
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

func (a *Auth) accessTokenPath(redirect, code string) string {
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

func (a *Auth) redirectServer(redirect chan<- requestAndError) (uri string, err error) {
	addr := a.RedirectHost
	if addr == "" {
		addr = DefaultRedirectHost
	}
	if port := a.RedirectPort; port != 0 {
		addr += ":" + strconv.Itoa(port)
	} else {
		addr += ":"
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return "", err
	}

	go func() {
		defer ln.Close()

		conn, err := ln.Accept()
		if err != nil {
			redirect <- requestAndError{nil, err}
			return
		}
		defer conn.Close()

		req, err := http.ReadRequest(bufio.NewReader(conn))
		redirect <- requestAndError{req, err}

		resp := http.Response{
			ProtoMajor: 1,
			ProtoMinor: 1,
			StatusCode: 200,
			Body: ioutil.NopCloser(strings.NewReader(
				"<script>window.close()</script>",
			)),
		}
		resp.Write(conn)
	}()

	return "http://" + ln.Addr().String(), nil
}

func browse(ctx context.Context, u string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", u)
	case "windows":
		cmd = exec.Command("rundll32", "u.dll,FileProtocolHandler", u)
	case "darwin":
		cmd = exec.Command("open", u)
	default:
		return fmt.Errorf("unsupported platform")
	}

	cmd.Start()

	ch := make(chan error, 1)
	go func() {
		ch <- cmd.Wait()
	}()
	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func codeFromParams(params url.Values) (string, error) {
	if errName := params.Get("error"); errName != "" {
		return "", fmt.Errorf(
			"bad redirect: %s: %s",
			errName, params.Get("error_description"),
		)
	}
	return params.Get("code"), nil
}

func waitRedirect(ctx context.Context, redirect <-chan requestAndError) (*http.Request, error) {
	select {
	case re := <-redirect:
		return re.req, re.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
