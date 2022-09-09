package httputils

import "net/http"

type JsonResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type JsonResponseError struct {
	Status int    `json:"-"`
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
}

func NewError(status, code int, msg string) error {
	return &JsonResponseError{
		Status: status,
		Code:   code,
		Msg:    msg,
	}
}

func (j *JsonResponseError) Error() string {
	return j.Msg
}

func NewBadRequestError(code int, err string) error {
	return &JsonResponseError{
		Status: http.StatusBadRequest,
		Code:   code,
		Msg:    err,
	}
}
