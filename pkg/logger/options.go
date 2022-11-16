package logger

import (
	"io"

	"go.uber.org/zap/zapcore"
)

type Option func(*Options)

type Options struct {
	w     io.Writer
	cores []zapcore.Core
	level string
}

type LokiConfig struct {
	Enable   bool
	Url      string
	Username string
	Password string
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
