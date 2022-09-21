package logger

import (
	"context"
	"net/url"

	"github.com/codfrm/cago/pkg/logger/loki"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Level string
	Debug bool
	Loki  *LokiConfig
}

func InitWithConfig(ctx context.Context, cfg *Config, opts ...Option) (*zap.Logger, error) {
	if cfg.Level != "" {
		opts = append(opts, Level(cfg.Level))
	}
	if cfg.Debug {
		opts = append(opts, Debug())
	}
	if cfg.Loki != nil {
		lokiOptions := make([]loki.Option, 0)
		u, err := url.Parse(cfg.Loki.Url)
		if err != nil {
			return nil, err
		}
		lokiOptions = append(lokiOptions, loki.WithLokiUrl(u))
		level := toLevel(cfg.Level)
		lokiOptions = append(lokiOptions, loki.WithLevelEnable(func(l zapcore.Level) bool {
			return l >= level
		}))
		lokiOptions = append(lokiOptions, loki.WithEnv())
		opts = append(opts, AppendCore(loki.NewLokiCore(ctx, lokiOptions...)))
	}
	return Init(opts...)
}

func Init(opt ...Option) (*zap.Logger, error) {
	options := &Options{}
	for _, o := range opt {
		o(options)
	}
	core := make([]zapcore.Core, 0, 1)
	level := toLevel(options.level)
	levelEnable := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= level
	})
	if options.w != nil {
		var encode zapcore.Encoder
		if options.debug {
			encode = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		} else {
			encode = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		}
		core = append(core, zapcore.NewCore(
			encode,
			zapcore.AddSync(options.w),
			levelEnable,
		))
	}
	if options.cores != nil {
		core = append(core, options.cores...)
	}
	logger := zap.New(zapcore.NewTee(core...))
	return logger, nil
}

func toLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	}
	return zap.InfoLevel
}
