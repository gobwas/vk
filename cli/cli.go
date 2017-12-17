package cli

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"unicode/utf8"

	"github.com/gobwas/vk"
)

func Authorize(ctx context.Context, app vk.App) (token *vk.AccessToken, err error) {
	redirect := make(chan requestAndError, 1)
	redirectPath, err := redirectServer(ctx, redirect)
	if err != nil {
		return nil, err
	}
	auth := app.AuthPathCode(redirectPath, vk.WithParam(
		"display", "page",
	))
	// Open a web browser to authorize an app.
	if err := Browse(ctx, auth); err != nil {
		return nil, err
	}
	req, err := waitRedirect(ctx, redirect)
	if err != nil {
		return nil, err
	}
	code, err := vk.CodeFromQuery(req.URL.Query())
	if err != nil {
		return nil, err
	}
	return app.Authorize(ctx, redirectPath, code)
}

func AuthorizeStandalone(ctx context.Context, app vk.App) (*vk.AccessToken, error) {
	auth := app.AuthPathToken(
		"https://oauth.vk.com/blank.html",
		vk.WithParam(
			"display", "page",
		),
	)
	// Open a web browser to authorize an app.
	if err := Browse(ctx, auth); err != nil {
		return nil, err
	}

	str, err := Ask(ctx, "Copy and paste url from browser: ")
	if err != nil {
		return nil, err
	}

	return vk.TokenFromURL(str)
}

func redirectServer(ctx context.Context, redirect chan<- requestAndError) (uri string, err error) {
	ln, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		return "", err
	}
	go http.Serve(ln, http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		defer ln.Close()
		redirect <- requestAndError{req, nil}
		io.Copy(rw, strings.NewReader(
			`<script>window.close()</script>`,
		))
	}))
	return "http://" + ln.Addr().String(), nil
}

func Browse(ctx context.Context, u string) error {
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

func Ask(ctx context.Context, question string) (string, error) {
	fmt.Fprint(os.Stderr, question)

	ch := make(chan stringAndError, 1)
	go func() {
		r := bufio.NewReader(os.Stdin)
		p, err := r.ReadBytes('\n')
		if n := len(p); n > 0 {
			p = p[:n-1]
		}
		ch <- stringAndError{string(p), err}
	}()

	select {
	case r := <-ch:
		return r.str, r.err
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func AskRune(ctx context.Context, question string) (rune, error) {
	// disable input buffering
	err := exec.Command("stty", "-f", "/dev/tty", "cbreak", "min", "1").Run()
	if err != nil {
		return 0, err
	}
	// do not display entered characters on the screen
	//exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	// restore the echoing state when exiting
	defer exec.Command("stty", "-F", "/dev/tty", "echo").Run()

	fmt.Fprint(os.Stderr, question)
	defer fmt.Fprint(os.Stderr, "\n")

	ch := make(chan runeAndError, 1)
	go func() {
		p := make([]byte, utf8.UTFMax)
		for i := range p {
			_, err := os.Stdin.Read(p[i : i+1])
			if err != nil {
				ch <- runeAndError{0, err}
				return
			}
			if utf8.FullRune(p) {
				break
			}
		}
		r, _ := utf8.DecodeRune(p)
		if r == utf8.RuneError {
			ch <- runeAndError{0, fmt.Errorf("invalid sequence")}
		} else {
			ch <- runeAndError{r, nil}
		}
	}()

	select {
	case r := <-ch:
		return r.r, r.err
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

func waitRequest(ctx context.Context, requests <-chan *http.Request) (*http.Request, error) {
	select {
	case req := <-requests:
		return req, nil
	case <-ctx.Done():
		return nil, ctx.Err()
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

type stringAndError struct {
	str string
	err error
}

type runeAndError struct {
	r   rune
	err error
}
