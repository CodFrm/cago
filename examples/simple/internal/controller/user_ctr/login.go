package user_ctr

import (
	"context"

	"github.com/gin-gonic/gin"

	api "github.com/codfrm/cago/examples/simple/internal/api/user"
	"github.com/codfrm/cago/examples/simple/internal/service/user_svc"
)

type Login struct {
}

func NewLogin() *Login {
	return &Login{}
}

// Register 注册
func (l *Login) Register(ctx context.Context, req *api.RegisterRequest) (*api.RegisterResponse, error) {
	return user_svc.Login().Register(ctx, req)
}

// Login 登录
func (l *Login) Login(ctx *gin.Context, req *api.LoginRequest) (*api.LoginResponse, error) {
	return user_svc.Login().Login(ctx, req)
}

// Logout 登出
func (l *Login) Logout(ctx *gin.Context, req *api.LogoutRequest) (*api.LogoutResponse, error) {
	return user_svc.Login().Logout(ctx, req)
}

func (l *Login) Middleware() gin.HandlerFunc {
	return user_svc.Login().Middleware()
}
