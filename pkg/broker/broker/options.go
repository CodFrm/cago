package broker

import "context"

type Option func(*Options)

type PublishOption func(options *PublishOptions)

type SubscribeOption func(options *SubscribeOptions)

type Options struct {
}

type PublishOptions struct {
	Context context.Context
}

type SubscribeOptions struct {
	Context    context.Context
	AutoAck    bool
	Group      string
	Retry      bool
	Concurrent int
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
		Context:    nil,
		AutoAck:    true,
		Group:      "",
		Retry:      false,
		Concurrent: 1,
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func Group(group string) SubscribeOption {
	return func(options *SubscribeOptions) {
		options.Group = group
	}
}

// NotAutoAck 不自动确认消息
func NotAutoAck() SubscribeOption {
	return func(options *SubscribeOptions) {
		options.AutoAck = false
	}
}

// Retry 产生错误时重试
func Retry() SubscribeOption {
	return func(options *SubscribeOptions) {
		options.Retry = true
	}
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

func WithConcurrent(concurrent int) SubscribeOption {
	return func(options *SubscribeOptions) {
		options.Concurrent = concurrent
	}
}
