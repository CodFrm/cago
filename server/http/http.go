package http

import (
	"context"
	"errors"

	"github.com/codfrm/cago"
	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/pkg/trace"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go.uber.org/zap"
)

type Config struct {
	Address []string `yaml:"address"`
}

type Callback func(r *Router) error

type server struct {
	ctx      context.Context
	cancel   context.CancelFunc
	callback Callback
}

// Http http服务组件,需要先注册logger组件
func Http(callback Callback) cago.ComponentCancel {
	return &server{
		callback: callback,
	}
}

func (h *server) Start(ctx context.Context, cfg *configs.Config) error {
	return h.StartCancel(ctx, nil, cfg)
}

func (h *server) StartCancel(
	ctx context.Context,
	cancel context.CancelFunc,
	cfg *configs.Config,
) error {
	config := &Config{}
	if err := cfg.Scan("http", config); err != nil {
		return err
	}
	l := logger.Default()
	r := gin.New()
	// 加入日志中间件
	r.Use(logger.Middleware(logger.Default()))
	if tp := trace.Default(); tp != nil {
		// 加入链路追踪中间件
		r.Use(trace.Middleware(cfg.AppName, tp))
	}
	if cfg.Env != configs.PROD {
		url := ginSwagger.URL("/swagger/doc.json")
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	}
	if err := h.callback(&Router{IRouter: r}); err != nil {
		return errors.New("failed to register http server")
	}
	// 启动http服务
	go func() {
		if len(config.Address) == 0 {
			config.Address = []string{"127.0.0.1:8080"}
		}
		if err := r.Run(config.Address...); err != nil {
			l.Error("failed to start http server", zap.Error(err))
			cancel()
		}
	}()
	return nil
}

func (h *server) CloseHandle() {
}
