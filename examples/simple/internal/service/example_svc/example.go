package example_svc

import (
	"context"
	"github.com/codfrm/cago/examples/simple/internal/service/user_svc"
	"github.com/codfrm/cago/pkg/utils"

	api "github.com/codfrm/cago/examples/simple/internal/api/example"
)

type ExampleSvc interface {
	// Login 需要登录的接口
	Login(ctx context.Context, req *api.LoginRequest) (*api.LoginResponse, error)
	// Ping ping
	Ping(ctx context.Context, req *api.PingRequest) (*api.PingResponse, error)
}

type exampleSvc struct {
}

var defaultExample = &exampleSvc{}

func Example() ExampleSvc {
	return defaultExample
}

// Login 需要登录的接口
func (e *exampleSvc) Login(ctx context.Context, req *api.LoginRequest) (*api.LoginResponse, error) {
	return &api.LoginResponse{
		Username: user_svc.Login().Get(ctx).Username,
	}, nil
}

// Ping ping
func (e *exampleSvc) Ping(ctx context.Context, req *api.PingRequest) (*api.PingResponse, error) {
	return &api.PingResponse{Pong: utils.RandString(6, utils.Mix)}, nil
}
