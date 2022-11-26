package user

import (
	"github.com/codfrm/cago/server/mux"
)

// LoginRequest 登录
type LoginRequest struct {
	mux.Meta `path:"/user/login" method:"POST"`
	// 用户名
	Username string `form:"username" binding:"required"`
}

type LoginResponse struct {
	// 用户名
	Username string `json:"username"`
}

// RegisterRequest 注册
type RegisterRequest struct {
	mux.Meta `path:"/user/register" method:"POST"`
	// 用户名
	Username string `form:"username" binding:"required"`
}

type RegisterResponse struct {
}

// LogoutRequest 登出
type LogoutRequest struct {
	mux.Meta `path:"/user/logout" method:"DELETE"`
}

type LogoutResponse struct {
}
