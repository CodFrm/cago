package loki

import (
	"net/url"
	"os"

	"go.uber.org/zap"
)

type Option func(*Options)

type Options struct {
	url      *url.URL
	level    zap.LevelEnablerFunc
	labels   []zap.Field
	username string
	password string
}

func BasicAuth(username string, password string) Option {
	return func(options *Options) {
		options.username = username
		options.password = password
	}
}

func WithLevelEnable(enab zap.LevelEnablerFunc) Option {
	return func(o *Options) {
		o.level = enab
	}
}

func WithLokiUrl(u *url.URL) Option {
	return func(o *Options) {
		o.url = u
	}
}

func AppendLabels(labels ...zap.Field) Option {
	return func(options *Options) {
		if options.labels == nil {
			options.labels = make([]zap.Field, 0)
		}
		options.labels = append(options.labels, labels...)
	}
}

func WithEnv() Option {
	return func(options *Options) {
		WithKubernetes()(options)
		WithHost()(options)
	}
}

func WithHost() Option {
	h, err := os.Hostname()
	if err != nil {
		return func(options *Options) {
		}
	}
	return AppendLabels(
		zap.String("hostname", h),
	)
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
