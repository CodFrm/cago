package user_ctr

import (
	"context"

	api "github.com/codfrm/cago/examples/simple/internal/api/user"
	"github.com/codfrm/cago/examples/simple/internal/service/user_svc"
)

type User struct {
}

func NewUser() *User {
	return &User{}
}

// Register 注册
func (u *User) Register(ctx context.Context, req *api.RegisterRequest) (*api.RegisterResponse, error) {
	return user_svc.User().Register(ctx, req)
}

// Logout 登出
func (u *User) Logout(ctx context.Context, req *api.LogoutRequest) (*api.LogoutResponse, error) {
	return user_svc.User().Logout(ctx, req)
}
