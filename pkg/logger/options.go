package logger

import "go.uber.org/zap"

type Option func(*Options)

type Options struct {
	level  string
	debug  bool
	loki   *LokiConfig
	fields []zap.Field
}

type LokiConfig struct {
	Level string
	Url   string
}

func WithLoki(config *LokiConfig) Option {
	return func(options *Options) {
		options.loki = config
	}
}

func Level(level string) Option {
	return func(options *Options) {
		options.level = level
	}
}

func Debug() Option {
	return func(options *Options) {
		options.debug = true
	}
}

func WithFields(fields ...zap.Field) Option {
	return func(options *Options) {
		options.fields = fields
	}
}
