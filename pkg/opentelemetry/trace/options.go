package trace

import (
	"go.opentelemetry.io/otel/attribute"
)

type Option func(*Options)

type Options struct {
	attrs []attribute.KeyValue
}

func AppendAttributes(attrs ...attribute.KeyValue) Option {
	return func(options *Options) {
		if options.attrs == nil {
			options.attrs = make([]attribute.KeyValue, 0)
		}
		options.attrs = append(options.attrs, attrs...)
	}
}
