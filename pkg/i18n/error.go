package i18n

import (
	"fmt"
	"net/http"
)

type Error struct {
	status int
	code   int
	args   []interface{}
}

func (e *Error) Error() string {
	return e.Msg(DefaultLang)
}

func (e *Error) Status() int {
	return e.status
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Msg(lang string) string {
	langMap, ok := langs[lang]
	if !ok {
		langMap = langs[DefaultLang]
	}
	return fmt.Sprintf(langMap[e.code], e.args...)
}

// NewErrorWithStatus 自定义错误
func NewErrorWithStatus(status int, code int, v ...interface{}) error {
	return &Error{
		status: status,
		code:   code,
		args:   v,
	}
}

// NewError 参数校验错误
func NewError(code int, v ...interface{}) error {
	return NewErrorWithStatus(http.StatusBadRequest, code, v...)
}

// NewUnauthorizedError 401未授权
func NewUnauthorizedError(code int, v ...interface{}) error {
	return NewErrorWithStatus(http.StatusUnauthorized, code, v...)
}

// NewForbiddenError 403禁止访问
func NewForbiddenError(code int, v ...interface{}) error {
	return NewErrorWithStatus(http.StatusForbidden, code, v...)
}

// NewNotFoundError 404资源不存在
func NewNotFoundError(code int, v ...interface{}) error {
	return NewErrorWithStatus(http.StatusNotFound, code, v...)
}

// NewInternalError 500构造内部错误
func NewInternalError(code int, v ...interface{}) error {
	return NewErrorWithStatus(http.StatusInternalServerError, code, v...)
}
