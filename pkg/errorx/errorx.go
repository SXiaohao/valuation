package errorx

import (
	"errors"
	"fmt"
)

type Error struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func New(code int, msg string) *Error {
	return &Error{
		Code: code,
		Msg:  msg,
	}
}

func (e Error) Error() string {
	return fmt.Sprintf("code: %d , msg: %s", e.Code, e.Msg)
}

func FromError(err error) *Error {
	if err == nil {
		return nil
	}
	if se := new(Error); errors.As(err, &se) {
		return se
	}
	return &Error{Code: 500, Msg: err.Error()}
}

func Bomb(code int, format string, a ...interface{}) {
	panic(New(code, fmt.Sprintf(format, a...)))
}

func Dangerous(v interface{}, code ...int) {
	if v == nil {
		return
	}

	c := 500
	if len(code) > 0 {
		c = code[0]
	}

	switch t := v.(type) {
	case string:
		if t != "" {
			panic(New(c, t))
		}
	case error:
		panic(New(c, t.Error()))
	case Error:
		panic(t)
	}
}
