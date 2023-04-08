package user

import "github.com/codfrm/cago/server/mux"

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
