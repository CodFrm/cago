package trace

import (
	"go.opentelemetry.io/otel/attribute"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

type Option func(*Options)

type Options struct {
	exp    tracesdk.SpanExporter
	sample float64
	attrs  []attribute.KeyValue
}

func WithExporter(exp tracesdk.SpanExporter) Option {
	return func(options *Options) {
		options.exp = exp
	}
}

func Sample(sample float64) Option {
	return func(options *Options) {
		options.sample = sample
	}
}

func AppendAttributes(attrs ...attribute.KeyValue) Option {
	return func(options *Options) {
		if options.attrs == nil {
			options.attrs = make([]attribute.KeyValue, 0)
		}
		options.attrs = append(options.attrs, attrs...)
	}
}
