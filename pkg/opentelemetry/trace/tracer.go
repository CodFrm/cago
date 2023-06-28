package trace

import (
	"context"
	"errors"

	exporter2 "github.com/codfrm/cago/pkg/opentelemetry/trace/exporter"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type ExporterType string

const (
	Jaeger  ExporterType = "jaeger"
	UpTrace ExporterType = "uptrace"
)

type Config struct {
	Endpoint string
	Type     ExporterType
	Username string
	Password string
	Dsn      string
	// Sample 采样率 0-1 其它数字为跟随父配置
	Sample float64
}

func NewWithConfig(ctx context.Context, cfg *Config, options ...Option) (trace.TracerProvider, error) {
	var exp tracesdk.SpanExporter
	var err error
	switch cfg.Type {
	case Jaeger:
		exp, err = exporter2.JaegerExporter(&exporter2.JaegerConfig{
			Endpoint: cfg.Endpoint,
			Username: cfg.Username,
			Password: cfg.Password,
		})
	case UpTrace:
		exp, err = exporter2.UpTraceExporter(&exporter2.UpTraceConfig{
			Dsn: cfg.Dsn,
		})
	default:
		return nil, errors.New("type not found")
	}
	if err != nil {
		return nil, err
	}
	options = append(options, Sample(cfg.Sample), WithExporter(exp))
	return New(options...)
}

func New(opt ...Option) (trace.TracerProvider, error) {
	options := &Options{}
	for _, o := range opt {
		o(options)
	}

	var sample tracesdk.Sampler
	if options.sample <= 0 {
		// 总是关闭
		sample = tracesdk.NeverSample()
	} else if options.sample < 1 {
		// 百分比采样，如果父开启了那么会开启
		sample = tracesdk.ParentBased(tracesdk.TraceIDRatioBased(options.sample))
	} else if options.sample == 1 {
		// 总是采样，如果父未开启那么不会开启
		sample = tracesdk.ParentBased(tracesdk.AlwaysSample())
	} else {
		// 总是关闭，如果父开启了那么会开启
		sample = tracesdk.ParentBased(tracesdk.NeverSample())
	}

	res, err := resource.New(context.Background(),
		//resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithHost(),
		//resource.WithContainer(),
		//resource.WithOS(),
		//resource.WithProcess(),
		resource.WithAttributes(
			options.attrs...,
		))
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(options.exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(res),
		tracesdk.WithSampler(sample),
	)
	return tp, nil
}
