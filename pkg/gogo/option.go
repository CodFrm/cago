package gogo

import "context"

type Option func(*Options)

type Options struct {
	ctx context.Context
}

// WithContext 设置上下文
func WithContext(ctx context.Context) Option {
	return func(o *Options) {
		o.ctx = ctx
	}
}
