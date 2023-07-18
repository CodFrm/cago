package cago

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/pkg/logger"
)

type Cago struct {
	ctx        context.Context
	cancel     context.CancelFunc
	cfg        *configs.Config
	components []Component
}

type CloseHandle func()

// New 初始化cago
func New(ctx context.Context, cfg *configs.Config) *Cago {
	ctx, cancel := context.WithCancel(ctx)
	cago := &Cago{
		ctx:    ctx,
		cancel: cancel,
		cfg:    cfg,
	}
	return cago
}

// Registry 注册组件
func (r *Cago) Registry(component Component) *Cago {
	err := component.Start(r.ctx, r.cfg)
	if err != nil {
		panic(err)
	}
	r.components = append(r.components, component)
	return r
}

// RegistryCancel 注册cancel组件
func (r *Cago) RegistryCancel(component ComponentCancel) *Cago {
	err := component.StartCancel(r.ctx, r.cancel, r.cfg)
	if err != nil {
		panic(errors.New("start component error: " + reflect.TypeOf(component).String() + " " + err.Error()))
	}
	r.components = append(r.components, component)
	return r
}

// Start 启动框架,在此之前组件已全部启动,此处只做停止等待
func (r *Cago) Start() error {
	quitSignal := make(chan os.Signal, 1)
	// 优雅启停
	signal.Notify(
		quitSignal,
		syscall.SIGINT, syscall.SIGTERM,
	)
	select {
	case <-quitSignal:
		r.cancel()
	case <-r.ctx.Done():
	}
	logger.Default().Info(r.cfg.AppName + " is stopping...")
	for _, v := range r.components {
		v.CloseHandle()
	}
	logger.Default().Info(r.cfg.AppName + " is stopped")
	return nil
}
