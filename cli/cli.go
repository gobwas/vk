package cli

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"

	"github.com/gobwas/vk"
)

func Authorize(ctx context.Context, app vk.App) (token *vk.AccessToken, err error) {
	redirect := make(chan requestAndError, 1)
	redirectPath, err := redirectServer(redirect)
	if err != nil {
		return nil, err
	}
	// Open a web browser to authorize an app.
	if err := browse(ctx, app.AuthorizePath(redirectPath)); err != nil {
		return nil, err
	}
	req, err := waitRedirect(ctx, redirect)
	if err != nil {
		return nil, err
	}
	code, err := vk.CodeFromParams(req.URL.Query())
	if err != nil {
		return nil, err
	}
	return app.Authorize(ctx, redirectPath, code)
}

func redirectServer(redirect chan<- requestAndError) (uri string, err error) {
	ln, err := net.Listen("tcp", "127.0.0.1:")
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

func waitRedirect(ctx context.Context, redirect <-chan requestAndError) (*http.Request, error) {
	select {
	case re := <-redirect:
		return re.req, re.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

type requestAndError struct {
	req *http.Request
	err error
}
