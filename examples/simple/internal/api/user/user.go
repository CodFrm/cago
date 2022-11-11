package user

import (
	"github.com/codfrm/cago/server/http"
)

// CreateUserRequest 创建用户
type CreateUserRequest struct {
	http.Route `path:"/user" method:"POST"`
	Name       string `form:"name" binding:"required"`
}

type CreateUserResponse struct {
}

type LoginRequest struct {
	http.Route `path:"/user/login" method:"POST"`
	Username   string `form:"username" binding:"required"`
}

type LoginResponse struct {
}
