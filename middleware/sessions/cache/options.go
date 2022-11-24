package cache

import (
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

type Options struct {
	prefix         string
	defaultMaxAge  int
	refreshTime    int
	sessionOptions *sessions.Options
	codecs         []securecookie.Codec
}

type Option func(options *Options)

func WithPrefix(prefix string) Option {
	return func(options *Options) {
		options.prefix = prefix
	}
}

// DefaultMaxAge 默认过期时间, 单位秒
func DefaultMaxAge(maxAge int) Option {
	return func(options *Options) {
		options.defaultMaxAge = maxAge
	}
}

func WithSessionOptions(options *sessions.Options) Option {
	return func(opts *Options) {
		opts.sessionOptions = options
	}
}

func WithCodecs(codecs ...securecookie.Codec) Option {
	return func(options *Options) {
		options.codecs = codecs
	}
}

// WithRefreshTime 设置刷新时间, 单位秒
func WithRefreshTime(refreshTime int) Option {
	return func(options *Options) {
		options.refreshTime = refreshTime
	}
}
