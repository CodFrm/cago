package cago

import (
	"context"

	"github.com/codfrm/cago/config"
)

type Cago struct {
	ctx        context.Context
	cancel     context.CancelFunc
	cfg        *config.Config
	components []Component
}

type CloseHandle func()

// New 初始化cago
func New(ctx context.Context, cfg *config.Config) *Cago {
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
		panic(err)
	}
	r.components = append(r.components, component)
	return r
}

// Start 启动框架,在此之前组件以全部启动,此处只做停止等待
func (r *Cago) Start() error {
	<-r.ctx.Done()
	for _, v := range r.components {
		v.CloseHandle()
	}
	return nil
}
