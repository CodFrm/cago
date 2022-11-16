package logger

import (
	"context"
	"io"
	"net/url"
	"os"

	"github.com/codfrm/cago/pkg/logger/loki"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	Level       string
	LogFile     LogFileConfig
	Loki        LokiConfig
	lokiOptions []loki.Option
	debug       bool
}

type LogFileConfig struct {
	Enable        bool
	Filename      string
	ErrorFilename string
}

func InitWithConfig(ctx context.Context, cfg *Config, opts ...Option) (*zap.Logger, error) {
	if cfg.Level != "" {
		opts = append(opts, Level(cfg.Level))
	}
	if cfg.LogFile.Enable {
		if cfg.LogFile.Filename != "" {
			opts = append(opts, AppendCore(NewFileCore(toLevel(cfg.Level), cfg.LogFile.Filename)))
		}
		if cfg.LogFile.ErrorFilename != "" {
			opts = append(opts, AppendCore(NewFileCore(zap.ErrorLevel, cfg.LogFile.ErrorFilename)))
		}
	}
	if cfg.debug {
		opts = append(opts, AppendCore(zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.Lock(os.Stdout),
			zapcore.DebugLevel,
		)))
	}
	if cfg.Loki.Enable {
		lokiOptions := cfg.lokiOptions
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
		if cfg.Loki.Username != "" {
			lokiOptions = append(lokiOptions, loki.BasicAuth(
				cfg.Loki.Username, cfg.Loki.Password,
			))
		}
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
		encodeConfig := zap.NewProductionEncoderConfig()
		encodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encode := zapcore.NewJSONEncoder(encodeConfig)
		core = append(core, zapcore.NewCore(
			encode,
			zapcore.AddSync(options.w),
			levelEnable,
		))
	}
	if options.cores != nil {
		core = append(core, options.cores...)
	}
	logger := zap.New(zapcore.NewTee(core...), zap.AddCaller())
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

func NewFileCore(level zapcore.Level, filename string) zapcore.Core {
	var w io.Writer = &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    2,
		MaxBackups: 10,
		MaxAge:     30,
		LocalTime:  true,
		Compress:   false,
	}
	levelEnable := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= level
	})
	encodeConfig := zap.NewProductionEncoderConfig()
	encodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encode := zapcore.NewJSONEncoder(encodeConfig)
	return zapcore.NewCore(
		encode,
		zapcore.AddSync(w),
		levelEnable,
	)
}
