package trace

import (
	"context"
	"errors"

	"github.com/codfrm/cago/configs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

type ExporterType string

const (
	Jaeger ExporterType = "jaeger"
)

type Config struct {
	Endpoint string
	Type     ExporterType
	Username string
	Password string
	// Sample 采样率 0-1 其它数字为跟随父配置
	Sample float64
}

func jaegerExporter(config *Config) (tracesdk.SpanExporter, error) {
	return jaeger.New(jaeger.WithCollectorEndpoint(
		jaeger.WithEndpoint(config.Endpoint),
		jaeger.WithUsername(config.Username),
		jaeger.WithPassword(config.Password),
	))
}

func InitWithConfig(ctx context.Context, config *configs.Config) (trace.TracerProvider, error) {
	cfg := &Config{}
	if err := config.Scan("trace", cfg); err != nil {
		return nil, err
	}
	var exp tracesdk.SpanExporter
	var err error
	switch cfg.Type {
	case Jaeger:
		exp, err = jaegerExporter(cfg)
	default:
		return nil, errors.New("type not found")
	}
	if err != nil {
		return nil, err
	}

	var sample tracesdk.Sampler
	if cfg.Sample <= 0 {
		// 总是关闭
		sample = tracesdk.NeverSample()
	} else if cfg.Sample < 1 {
		// 百分比采样，如果父开启了那么会开启
		sample = tracesdk.ParentBased(tracesdk.TraceIDRatioBased(cfg.Sample))
	} else if cfg.Sample == 1 {
		// 总是采样，如果父未开启那么不会开启
		sample = tracesdk.ParentBased(tracesdk.AlwaysSample())
	} else {
		// 总是关闭，如果父开启了那么会开启
		sample = tracesdk.ParentBased(tracesdk.NeverSample())
	}

	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(config.AppName),
			attribute.String("environment", string(config.Env)),
		)),
		tracesdk.WithSampler(sample),
	)
	return tp, nil
}
