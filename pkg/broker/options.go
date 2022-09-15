package broker

import "go.opentelemetry.io/otel/trace"

type Option func(options *Options)

type Options struct {
	traceProvider trace.TracerProvider
}

func WithTraceProvider(t trace.TracerProvider) Option {
	return func(options *Options) {
		options.traceProvider = t
	}
}
