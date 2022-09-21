package logger

import (
	"context"
	"os"

	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/pkg/logger/loki"
	"go.uber.org/zap"
)

var logger *zap.Logger

// Logger 日志组件,核心组件,必须注册
func Logger(ctx context.Context, config *configs.Config) error {
	cfg := &Config{}
	if err := config.Scan("logger", cfg); err != nil {
		return err
	}
	l, err := InitWithConfig(ctx, cfg, WithWriter(os.Stdout),
		WithLokiOptions(loki.AppendLabels(zap.String("app", config.AppName))))
	if err != nil {
		return err
	}
	logger = l
	return nil
}

func SetLogger(l *zap.Logger) {
	logger = l
}

func Default() *zap.Logger {
	return logger
}
