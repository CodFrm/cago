package httputils

import (
	"net/http"
)

// JSONResponse 返回json数据
type JSONResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Error 错误返回，包含了状态码、错误码、错误信息、请求ID
// 请求id为trace链路追踪的id，需要配置trace组件才会有
type Error struct {
	Status    int    `json:"-"`
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	RequestID string `json:"request_id,omitempty"`
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
