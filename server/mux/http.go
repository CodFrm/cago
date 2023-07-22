package mux

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/codfrm/cago/pkg/gogo"

	"github.com/codfrm/cago"
	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/pkg/opentelemetry/metric"
	"github.com/codfrm/cago/pkg/opentelemetry/trace"
	"github.com/codfrm/cago/pkg/utils/validator"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

type Config struct {
	Address []string `yaml:"address"`
}

type Callback func(r *Router) error

type server struct {
	//ctx context.Context
	//cancel   context.CancelFunc
	callback Callback
}

// HTTP http服务组件,需要先注册logger组件
func HTTP(callback Callback) cago.ComponentCancel {
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
	err := cfg.Scan("http", config)
	if err != nil {
		return err
	}
	l := logger.Default()
	var r *gin.Engine
	if cfg.Debug {
		gin.SetMode(gin.DebugMode)
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
		r.Use(Recover())
	}
	binding.Validator, err = validator.NewValidator()
	if err != nil {
		return err
	}
	// ginContext支持fallback
	r.ContextWithFallback = true
	// 加入日志中间件
	r.Use(logger.Middleware(logger.Default()))
	// 加入metrics中间件
	if metric.Default() != nil {
		r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	}

	if tp := trace.Default(); tp != nil {
		// 加入链路追踪中间件
		r.Use(trace.Middleware(cfg.AppName, tp))
	}
	if cfg.Env != configs.PROD {
		url := ginSwagger.URL("/swagger/doc.json")
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	}
	if err := h.callback(&Router{IRouter: r}); err != nil {
		return errors.New("failed to register http server: " + err.Error())
	}

	if len(config.Address) == 0 {
		config.Address = []string{"127.0.0.1:80"}
	}
	srv := &http.Server{
		Addr:              config.Address[0],
		Handler:           r.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	_ = gogo.Go(func(ctx context.Context) error {
		<-ctx.Done()
		l.Info("http server closing...")
		if err := srv.Shutdown(context.Background()); err != nil {
			l.Error("failed to close http server", zap.Error(err))
			return err
		}
		l.Info("http server closed")
		return nil
	}, gogo.WithContext(ctx))
	// 启动http服务
	_ = gogo.Go(func(ctx context.Context) error {
		defer cancel()
		if err := srv.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				l.Info("http server closed")
				return nil
			}
			l.Error("failed to start http server", zap.Error(err))
			return err
		}
		return nil
	})
	return nil
}

func (h *server) CloseHandle() {

}
