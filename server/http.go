package server

import (
	"context"
	"errors"

	"github.com/codfrm/cago"
	"github.com/codfrm/cago/config"
	"github.com/codfrm/cago/mux"
	"github.com/codfrm/cago/pkg/logger"
	"go.uber.org/zap"
)

type http struct {
	ctx      context.Context
	cancel   context.CancelFunc
	callback func(r *mux.RouterGroup) error
}

// Http http服务组件,需要先注册logger组件
func Http(callback func(r *mux.RouterGroup) error) cago.ComponentCancel {
	return &http{
		callback: callback,
	}
}

func (h *http) Start(ctx context.Context, cfg *config.Config) error {
	panic("not implement")
}

func (h *http) StartCancel(ctx context.Context, cancel context.CancelFunc, cfg *config.Config) error {
	l := logger.Ctx(cago.Background()).With(
		zap.String("app", cfg.AppName),
	)
	mux := mux.New(l)
	if err := h.callback(mux.Group()); err != nil {
		return errors.New("failed to register http")
	}
	// 启动http服务
	go func() {
		if err := mux.Run(":8080"); err != nil {
			l.Error("failed to start http", zap.Error(err))
			cancel()
		}
	}()
	return nil
}

func (h *http) CloseHandle() {
}
