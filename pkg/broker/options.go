package broker

import (
	broker2 "github.com/codfrm/cago/pkg/broker/broker"
	"go.opentelemetry.io/otel/trace"
)

type Option func(options *Options)

type Options struct {
	tracer trace.Tracer
	broker broker2.Broker
}

func WithTracer(t trace.Tracer) Option {
	return func(options *Options) {
		options.tracer = t
	}
}

func WithBroker(b broker2.Broker) Option {
	return func(options *Options) {
		options.broker = b
	}
}
