package controller

import (
	"context"
	"errors"

	"github.com/codfrm/cago/examples/simple/internal/api"
)

type User struct {
}

func NewUser() User {
	return User{}
}

// CreateUser 创建用户
func (u *User) CreateUser(ctx context.Context, req *api.CreateUserRequest) (*api.CreateUserResponse, error) {

	return nil, errors.New("not implement")
}
