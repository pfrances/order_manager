package domain

import (
	"errors"
	"fmt"
)

var (
	EUNKNOWN  = "EUNKNOWN"
	ENOTFOUND = "ENOTFOUND"
	EINVALID  = "EINVALID"
)

type Error struct {
	code    string
	message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.code, e.message)
}

func ErrorCode(err error) string {
	var e *Error
	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.code
	}
	return EUNKNOWN
}

func Errorf(code string, format string, args ...interface{}) *Error {
	return &Error{
		code:    code,
		message: fmt.Sprintf(format, args...),
	}
}
