package service

import (
	"context"

	"github.com/codfrm/cago/examples/simple/internal/api"
)

type IUser interface {
	// Create 创建用户
	Create(ctx context.Context, req *api.CreateUserRequest) (*api.CreateUserResponse, error)
}

type user struct {
}

var defaultUser = &user{}

func User() IUser {
	return defaultUser
}

// Create 创建用户
func (u *user) Create(ctx context.Context, req *api.CreateUserRequest) (*api.CreateUserResponse, error) {
	return nil, nil
}
