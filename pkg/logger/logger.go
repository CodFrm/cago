package logger

import (
	"context"
	"net/url"
	"os"

	"github.com/codfrm/gocat/pkg/logger/loki"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger(ctx context.Context, opt ...Option) (*zap.Logger, error) {
	options := &Options{}
	for _, o := range opt {
		o(options)
	}
	core := make([]zapcore.Core, 1)
	level := toLevel(options.level)
	levelEnable := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= level
	})
	if options.debug {
		core = append(core, zapcore.NewCore(
			zapcore.NewConsoleEncoder(zapcore.EncoderConfig{}),
			zapcore.AddSync(os.Stdout),
			levelEnable,
		))
	} else {
		core = append(core, zapcore.NewCore(
			zapcore.NewJSONEncoder(zapcore.EncoderConfig{}),
			zapcore.AddSync(os.Stdout),
			levelEnable,
		))
	}
	if options.loki != nil {
		u, err := url.Parse(options.loki.Url)
		if err != nil {
			return nil, err
		}
		level := toLevel(options.loki.Level)
		lokiCore, err := loki.NewLokiCore(ctx, u, zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= level
		}))
		if err != nil {
			return nil, err
		}
		core = append(core, lokiCore)
	}

	logger := zap.New(zapcore.NewTee(core...))
	Logger = logger
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
