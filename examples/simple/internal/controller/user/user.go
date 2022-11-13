package user

import (
	"context"

	api "github.com/codfrm/cago/examples/simple/internal/api/user"
	service "github.com/codfrm/cago/examples/simple/internal/service/user"
)

type User struct {
}

func NewUser() User {
	return User{}
}

// Login 登录
func (u *User) Login(ctx context.Context, req *api.LoginRequest) (*api.LoginResponse, error) {
	return service.User().Login(ctx, req)
}

// Register 注册
func (u *User) Register(ctx context.Context, req *api.RegisterRequest) (*api.RegisterResponse, error) {
	return service.User().Register(ctx, req)
}

// Logout 登出
func (u *User) Logout(ctx context.Context, req *api.LogoutRequest) (*api.LogoutResponse, error) {
	return service.User().Logout(ctx, req)
}
