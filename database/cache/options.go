package cache

import "time"

type Option func(*Options)

type Options struct {
	expiration time.Duration
	depend     Depend
}

func NewOptions(opts ...Option) *Options {
	options := &Options{}
	for _, v := range opts {
		v(options)
	}
	return options
}

func Expiration(t time.Duration) Option {
	return func(options *Options) {
		options.expiration = t
	}
}

func WithDepend(depend Depend) Option {
	return func(options *Options) {
		options.depend = depend
	}
}
