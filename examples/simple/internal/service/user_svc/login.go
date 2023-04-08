package user_svc

import (
	"context"

	api "github.com/codfrm/cago/examples/simple/internal/api/user"
)

type LoginSvc interface {
	// Login 登录
	Login(ctx context.Context, req *api.LoginRequest) (*api.LoginResponse, error)
}

type loginSvc struct {
}

var defaultLogin = &loginSvc{}

func Login() LoginSvc {
	return defaultLogin
}

// Login 登录
func (l *loginSvc) Login(ctx context.Context, req *api.LoginRequest) (*api.LoginResponse, error) {
	return nil, nil
}
