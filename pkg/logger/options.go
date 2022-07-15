package logger

import "go.uber.org/zap"

type Option func(*Options)

type Options struct {
	level  string
	debug  bool
	loki   *LokiConfig
	labels []zap.Field
}

type LokiConfig struct {
	Level string `yaml:"level"`
	Url   string `yaml:"url"`
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

func WithLabels(labels ...zap.Field) Option {
	return func(options *Options) {
		options.labels = labels
	}
}
