package logger

import (
	"os"

	"go.uber.org/zap"
)

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

func AppendLabels(labels ...zap.Field) Option {
	return func(options *Options) {
		if options.labels == nil {
			options.labels = make([]zap.Field, 0)
		}
		for _, v := range labels {
			options.labels = append(options.labels, v)
		}
	}
}

func WithKubernetes() Option {
	namespace := os.Getenv("KUBERNETES_NAMESPACE")
	if namespace == "" {
		return func(options *Options) {}
	}
	return AppendLabels(
		zap.String("namespace", namespace),
		zap.String("pod", os.Getenv("KUBERNETES_POD_NAME")),
		zap.String("node", os.Getenv("KUBERNETES_NODE_NAME")),
	)
}
