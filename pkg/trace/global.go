package trace

import (
	"context"

	"github.com/codfrm/cago/configs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var tracerProvider trace.TracerProvider

// Trace 链路追踪组件
func Trace(ctx context.Context, config *configs.Config) error {
	tp, err := InitWithConfig(ctx, config)
	if err != nil {
		return err
	}
	tracerProvider = tp
	otel.SetTracerProvider(tp)
	return nil
}

func SetTraceProvider(tp trace.TracerProvider) {
	tracerProvider = tp
}

func Default() trace.TracerProvider {
	return tracerProvider
}
