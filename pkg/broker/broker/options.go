package broker

import "context"

type Option func(*Options)

type PublishOption func(options *PublishOptions)

type SubscribeOption func(options *SubscribeOptions)

type Options struct {
	DefaultGroup string
	TopicPrefix  string
}

type PublishOptions struct {
	Context context.Context
}

type SubscribeOptions struct {
	Context context.Context
	AutoAck bool
	Group   string
}

func NewOptions(opts ...Option) Options {
	opt := Options{}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func NewPublishOptions(opts ...PublishOption) PublishOptions {
	opt := PublishOptions{}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func NewSubscribeOptions(opts ...SubscribeOption) SubscribeOptions {
	opt := SubscribeOptions{
		AutoAck: true,
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func WithPublishContext(ctx context.Context) PublishOption {
	return func(options *PublishOptions) {
		options.Context = ctx
	}
}

func WithSubscribeContext(ctx context.Context) SubscribeOption {
	return func(options *SubscribeOptions) {
		options.Context = ctx
	}
}
