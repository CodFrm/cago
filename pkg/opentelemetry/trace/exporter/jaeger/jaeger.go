package jaeger

import (
	"context"

	"github.com/codfrm/cago/pkg/opentelemetry/trace"
	"go.opentelemetry.io/otel/exporters/jaeger"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

type Config struct {
	Endpoint string
	Username string
	Password string
}

func init() {
	trace.RegisterExporter("jaeger", func(ctx context.Context, config *trace.Config) (tracesdk.SpanExporter, error) {
		return Exporter(&Config{
			Endpoint: config.Endpoint,
			Username: config.Username,
			Password: config.Password,
		})
	})
}

func Exporter(config *Config) (tracesdk.SpanExporter, error) {
	return jaeger.New(jaeger.WithCollectorEndpoint(
		jaeger.WithEndpoint(config.Endpoint),
		jaeger.WithUsername(config.Username),
		jaeger.WithPassword(config.Password),
	))
}
