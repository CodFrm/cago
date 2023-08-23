package trace

import (
	"context"
	"github.com/codfrm/cago/server/mux"
	"github.com/gin-gonic/gin"

	tracesdk "go.opentelemetry.io/otel/sdk/trace"

	"github.com/codfrm/cago/configs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

func init() {
	mux.RegisterMiddleware(func(cfg *configs.Config, r *gin.Engine) error {
		if tp := Default(); tp != nil {
			// 加入链路追踪中间件
			r.Use(Middleware(cfg.AppName, tp))
		}
		return nil
	})
}

type NewExporterFunc func(ctx context.Context, config *Config) (tracesdk.SpanExporter, error)

var (
	exporters = make(map[string]NewExporterFunc)
)

func RegisterExporter(name string, f NewExporterFunc) {
	exporters[name] = f
}

var tracerProvider trace.TracerProvider

// Trace 链路追踪组件,尽早注册,其他组件会判断是否存在,再进行使用
func Trace(ctx context.Context, config *configs.Config) error {
	cfg := &Config{}
	if err := config.Scan("trace", cfg); err != nil {
		return err
	}
	tp, err := NewWithConfig(ctx, cfg, AppendAttributes(
		semconv.ServiceNameKey.String(config.AppName),
		semconv.ServiceVersionKey.String(config.Version),
		semconv.DeploymentEnvironmentKey.String(string(config.Env)),
	))
	if err != nil {
		return err
	}
	tracerProvider = tp
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		//propagation.Baggage{},
	))
	return nil
}

func Default() trace.TracerProvider {
	return tracerProvider
}
