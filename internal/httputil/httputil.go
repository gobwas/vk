package httputil

import (
	"fmt"
	"net/http"
)

func CheckResponseStatus(resp *http.Response) error {
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return fmt.Errorf(
		"bad status code: %d %q",
		resp.StatusCode, resp.Status,
	)
}
