package mux

import "go.opentelemetry.io/otel/trace"

type Option func(*Options)

type Options struct {
	serviceName    string
	tracerProvider trace.TracerProvider
}

func ServiceName(serviceName string) Option {
	return func(options *Options) {
		options.serviceName = serviceName
	}
}

// WithTracerProvider 开启链路追踪
func WithTracerProvider(tracer trace.TracerProvider) Option {
	return func(options *Options) {
		options.tracerProvider = tracer
	}
}
