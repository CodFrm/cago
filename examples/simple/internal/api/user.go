package api

import (
	"github.com/codfrm/cago/server/http"
)

// CreateUserRequest 创建用户
type CreateUserRequest struct {
	http.Route `path:"/user" method:"POST"`
	// Username 用户名
	Username string `json:"username" validate:"required"`
}

type CreateUserResponse struct {
}

type UserInfoRequest struct {
	http.Route `path:"/user/:uid" method:"GET"`
	Uid        int64 `in:"path"`
}

type UserInfoResponse struct {
	Username string `json:"username"`
}
