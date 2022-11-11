package user

import (
	"context"

	api "github.com/codfrm/cago/examples/simple/internal/api/user"
)

type IUser interface {
	// CreateUser 创建用户
	CreateUser(ctx context.Context, req *api.CreateUserRequest) (*api.CreateUserResponse, error)
	// Login TODO
	Login(ctx context.Context, req *api.LoginRequest) (*api.LoginResponse, error)
}

type user struct {
}

var defaultUser = &user{}

func User() IUser {
	return defaultUser
}

// CreateUser 创建用户
func (u *user) CreateUser(ctx context.Context, req *api.CreateUserRequest) (*api.CreateUserResponse, error) {
	return nil, nil
}

// Login TODO
func (u *user) Login(ctx context.Context, req *api.LoginRequest) (*api.LoginResponse, error) {
	return nil, nil
}
