package user

import (
	"github.com/codfrm/cago/server/http"
)

// LoginRequest 登录
type LoginRequest struct {
	http.Route `path:"/user/login" method:"POST"`
	// 用户名
	Username string `form:"username" binding:"required"`
}

type LoginResponse struct {
	// 用户名
	Username string `json:"username"`
}

// RegisterRequest 注册
type RegisterRequest struct {
	http.Route `path:"/user/register" method:"POST"`
	// 用户名
	Username string `form:"username" binding:"required"`
}

type RegisterResponse struct {
}

// LogoutRequest 登出
type LogoutRequest struct {
	http.Route `path:"/user/logout" method:"DELETE"`
}

type LogoutResponse struct {
}
