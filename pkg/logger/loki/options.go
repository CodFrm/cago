package loki

import "go.uber.org/zap/zapcore"

type Option func(*Options)

type Options struct {
	sync zapcore.WriteSyncer
}

func WithSync(sync zapcore.WriteSyncer) Option {
	return func(options *Options) {
		options.sync = sync
	}
}
