package exporter

import (
	"go.opentelemetry.io/otel/exporters/jaeger"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

type JaegerConfig struct {
	Endpoint string
	Username string
	Password string
}

func JaegerExporter(config *JaegerConfig) (tracesdk.SpanExporter, error) {
	return jaeger.New(jaeger.WithCollectorEndpoint(
		jaeger.WithEndpoint(config.Endpoint),
		jaeger.WithUsername(config.Username),
		jaeger.WithPassword(config.Password),
	))
}
