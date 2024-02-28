package logger

import (
	"context"
	"os"

	"github.com/codfrm/cago/configs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerContextKeyType int

const loggerKey loggerContextKeyType = iota

type InitLogger func(ctx context.Context, config *configs.Config, loggerConfig *Config) ([]Option, error)

var (
	logger     *zap.Logger
	initLogger = make([]InitLogger, 0)
)

func RegistryInitLogger(f InitLogger) {
	initLogger = append(initLogger, f)
}

// Logger 日志组件,核心组件,必须提前注册
func Logger(ctx context.Context, config *configs.Config) error {
	cfg := &Config{}
	if err := config.Scan(ctx, "logger", cfg); err != nil {
		return err
	}
	opts := make([]Option, 0)
	if cfg.Level != "" {
		opts = append(opts, Level(cfg.Level))
	}
	if cfg.LogFile.Enable {
		if cfg.LogFile.Filename != "" {
			opts = append(opts, AppendCore(NewFileCore(ToLevel(cfg.Level), cfg.LogFile.Filename)))
		}
		if cfg.LogFile.ErrorFilename != "" {
			opts = append(opts, AppendCore(NewFileCore(zap.ErrorLevel, cfg.LogFile.ErrorFilename)))
		}
	}
	if config.Debug {
		opts = append(opts, AppendCore(zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.Lock(os.Stdout),
			zapcore.DebugLevel,
		)))
	}
	for _, f := range initLogger {
		o, err := f(ctx, config, cfg)
		if err != nil {
			return err
		}
		opts = append(opts, o...)
	}
	l, err := New(opts...)
	if err != nil {
		return err
	}
	logger = l
	return nil
}

func SetLogger(l *zap.Logger) {
	logger = l
}

// Default 默认日志,尽量不要使用,会丢失上下文信息
func Default() *zap.Logger {
	return logger
}

func Ctx(ctx context.Context) *zap.Logger {
	log, ok := ctx.Value(loggerKey).(*zap.Logger)
	if !ok {
		return logger
	}
	return log
}

func ContextWithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}
