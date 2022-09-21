package logger

import (
	"io"

	"github.com/codfrm/cago/pkg/logger/loki"
	"go.uber.org/zap/zapcore"
)

type Option func(*Options)

type Options struct {
	w           io.Writer
	cores       []zapcore.Core
	level       string
	debug       bool
	lokiOptions []loki.Option
}

type LokiConfig struct {
	Url string
}

func WithLokiOptions(opt ...loki.Option) Option {
	return func(o *Options) {
		o.lokiOptions = append(o.lokiOptions, opt...)
	}
}

func WithWriter(w io.Writer) Option {
	return func(o *Options) {
		o.w = w
	}
}

func AppendCore(core ...zapcore.Core) Option {
	return func(o *Options) {
		if o.cores == nil {
			o.cores = make([]zapcore.Core, 0)
		}
		o.cores = append(o.cores, core...)
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
