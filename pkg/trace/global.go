package trace

import (
	"context"

	"github.com/codfrm/cago/configs"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

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
	return nil
}

func SetTracerProvider(tp trace.TracerProvider) {
	tracerProvider = tp
}

func Default() trace.TracerProvider {
	return tracerProvider
}
