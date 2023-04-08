package user_svc

import (
	"context"

	api "github.com/codfrm/cago/examples/simple/internal/api/user"
)

type UserSvc interface {
	// Register 注册
	Register(ctx context.Context, req *api.RegisterRequest) (*api.RegisterResponse, error)
	// Logout 登出
	Logout(ctx context.Context, req *api.LogoutRequest) (*api.LogoutResponse, error)
}

type userSvc struct {
}

var defaultUser = &userSvc{}

func User() UserSvc {
	return defaultUser
}

// Register 注册
func (u *userSvc) Register(ctx context.Context, req *api.RegisterRequest) (*api.RegisterResponse, error) {
	return nil, nil
}

// Logout 登出
func (u *userSvc) Logout(ctx context.Context, req *api.LogoutRequest) (*api.LogoutResponse, error) {
	return nil, nil
}
