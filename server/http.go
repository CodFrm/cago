package server

import (
	"context"
	"errors"

	"github.com/codfrm/cago"
	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/mux"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/pkg/trace"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go.uber.org/zap"
)

type HttpConfig struct {
	Address []string `yaml:"address"`
}

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

func (h *http) Start(ctx context.Context, cfg *configs.Config) error {
	return h.StartCancel(ctx, nil, cfg)
}

func (h *http) StartCancel(
	ctx context.Context,
	cancel context.CancelFunc,
	cfg *configs.Config,
) error {
	config := &HttpConfig{}
	if err := cfg.Scan("http", config); err != nil {
		return err
	}
	l := logger.Default()
	opts := []mux.Option{
		mux.ServiceName(cfg.AppName),
	}
	if tp := trace.Default(); tp != nil {
		opts = append(opts, mux.WithTracerProvider(tp))
	}
	r := mux.New(l, opts...)
	group := r.Group()
	if cfg.Env != configs.PROD {
		url := ginSwagger.URL("/swagger/doc.json")
		group.GET("/swagger/*any", func(c *mux.WebContext) {
			ginSwagger.WrapHandler(swaggerFiles.Handler, url)(c.Context)
		})
	}
	if err := h.callback(group); err != nil {
		return errors.New("failed to register http")
	}
	// 启动http服务
	go func() {
		if len(config.Address) == 0 {
			config.Address = []string{"127.0.0.1:8080"}
		}
		if err := r.Run(config.Address...); err != nil {
			l.Error("failed to start http", zap.Error(err))
			cancel()
		}
	}()
	return nil
}

func (h *http) CloseHandle() {
}
