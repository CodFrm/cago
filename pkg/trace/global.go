package trace

import (
	"context"

	"github.com/codfrm/cago/configs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

var tracerProvider trace.TracerProvider

// Trace 链路追踪组件
func Trace(ctx context.Context, config *configs.Config) error {
	cfg := &Config{}
	if err := config.Scan("trace", cfg); err != nil {
		return err
	}
	tp, err := InitWithConfig(ctx, cfg, AppendAttributes(
		semconv.ServiceNameKey.String(config.AppName),
		attribute.String("environment", string(config.Env)),
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
