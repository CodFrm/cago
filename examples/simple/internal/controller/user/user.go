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

// CreateUser 创建用户
func (u *User) CreateUser(ctx context.Context, req *api.CreateUserRequest) (*api.CreateUserResponse, error) {
	return service.User().CreateUser(ctx, req)
}

// Login TODO
func (u *User) Login(ctx context.Context, req *api.LoginRequest) (*api.LoginResponse, error) {
	return service.User().Login(ctx, req)
}
