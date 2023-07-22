package httputils

import (
	"net/http"
)

type JSONResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type Error struct {
	Status int    `json:"-"`
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
}

func NewError(status, code int, msg string) error {
	return &Error{
		Status: status,
		Code:   code,
		Msg:    msg,
	}
}

func (j *Error) Error() string {
	return j.Msg
}

func NewBadRequestError(code int, err string) error {
	return NewError(http.StatusBadRequest, code, err)
}

func NewUnauthorizedError(code int, err string) error {
	return NewError(http.StatusUnauthorized, code, err)
}

func NewForbiddenError(code int, err string) error {
	return NewError(http.StatusForbidden, code, err)
}

func NewNotFoundError(code int, err string) error {
	return NewError(http.StatusNotFound, code, err)
}

func NewInternalServerError(code int, err string) error {
	return NewError(http.StatusInternalServerError, code, err)
}
