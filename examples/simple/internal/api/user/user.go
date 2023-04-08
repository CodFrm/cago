package user

import (
	"github.com/codfrm/cago/server/mux"
)

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
