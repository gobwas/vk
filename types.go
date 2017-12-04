package vk

import (
	"fmt"
	"strings"

	"github.com/mailru/easyjson"
)

//go:generate easyjson -all

type Response struct {
	Body easyjson.RawMessage `json:"response"`
	Err  Error               `json:"error"`
}

func (r Response) Error() *Error {
	if r.Err.Code == 0 {
		return nil
	}
	return &r.Err
}

type Error struct {
	Code   ErrorCode      `json:"error_code"`
	Msg    string         `json:"error_msg"`
	Params []RequestParam `json:"request_params"`
}

func (e Error) Error() string {
	return fmt.Sprintf(
		"%s (%d)",
		strings.ToLower(e.Msg), e.Code,
	)
}

func (e Error) Temporary() bool {
	switch e.Code {
	case ErrRateLimitExceeded, ErrTooManyActions:
		return true
	default:
		return false
	}
}

type RequestParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

//easyjson:json
type rawAccess struct {
	Token   string `json:"access_token"`
	Expires int    `json:"expires_in"`
	UserID  int    `json:"user_id"`
}
