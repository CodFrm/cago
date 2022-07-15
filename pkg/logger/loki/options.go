package loki

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Option func(*Options)

type Options struct {
	sync   zapcore.WriteSyncer
	labels []zap.Field
}

func WithSync(sync zapcore.WriteSyncer) Option {
	return func(options *Options) {
		options.sync = sync
	}
}

func WithLabels(labels ...zap.Field) Option {
	return func(options *Options) {
		options.labels = labels
	}
}
