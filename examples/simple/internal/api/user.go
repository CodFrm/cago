package api

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
