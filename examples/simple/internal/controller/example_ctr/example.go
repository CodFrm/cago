package example_ctr

import (
	"context"

	api "github.com/codfrm/cago/examples/simple/internal/api/example"
	"github.com/codfrm/cago/examples/simple/internal/service/example_svc"
)

type Example struct {
}

func NewExample() *Example {
	return &Example{}
}

// Login 需要登录的接口
func (e *Example) Login(ctx context.Context, req *api.LoginRequest) (*api.LoginResponse, error) {
	return example_svc.Example().Login(ctx, req)
}

// Ping ping
func (e *Example) Ping(ctx context.Context, req *api.PingRequest) (*api.PingResponse, error) {
	return example_svc.Example().Ping(ctx, req)
}
