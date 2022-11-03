package controller

import (
	"context"
	"errors"

	"github.com/codfrm/cago/examples/simple/internal/api"
)

type User struct {
}

// CreateUser 创建用户
func (u *User) CreateUser(ctx context.Context, req *api.CreateUserRequest) (*api.CreateUserResponse, error) {

	return nil, errors.New("not implement")
}

// UserInfo  在api中没有找到注释
func (u *User) UserInfo(ctx context.Context, req *api.UserInfoRequest) (*api.UserInfoResponse, error) {

	return nil, errors.New("not implement")
}
