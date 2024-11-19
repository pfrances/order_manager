package domain

import (
	"context"
	"errors"
	"fmt"
)

var (
	EUNKNOWN  = "EUNKNOWN"
	ENOTFOUND = "ENOTFOUND"
	EINVALID  = "EINVALID"
	ECANCELED = "ECANCELED"
)

type Error struct {
	code    string
	message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.code, e.message)
}

func ErrorCode(err error) string {
	if err == nil {
		return ""
	}

	var e *Error
	if errors.As(err, &e) {
		return e.code
	}

	if errors.Is(err, context.Canceled) {
		return ECANCELED
	}
	return EUNKNOWN
}

func Errorf(code string, format string, args ...interface{}) *Error {
	return &Error{
		code:    code,
		message: fmt.Sprintf(format, args...),
	}
}
