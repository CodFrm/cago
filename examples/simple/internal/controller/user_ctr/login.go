package user_ctr

import (
	"context"

	api "github.com/codfrm/cago/examples/simple/internal/api/user"
	"github.com/codfrm/cago/examples/simple/internal/service/user_svc"
)

type Login struct {
}

func NewLogin() *Login {
	return &Login{}
}

// Login 登录
func (l *Login) Login(ctx context.Context, req *api.LoginRequest) (*api.LoginResponse, error) {
	return user_svc.Login().Login(ctx, req)
}
