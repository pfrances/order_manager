package main

import "fmt"

var (
	ENOTFOUND = "ENOTFOUND"
)

type Error struct {
	code    string
	message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.code, e.message)
}

func Errorf(code string, format string, args ...interface{}) *Error {
	return &Error{
		code:    code,
		message: fmt.Sprintf(format, args...),
	}
}
