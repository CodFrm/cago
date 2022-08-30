package logger

import (
	"context"
	"os"

	"github.com/codfrm/cago/configs"
	"go.uber.org/zap"
)

var logger *zap.Logger

// Logger 日志组件,核心组件,必须注册
func Logger(ctx context.Context, config *configs.Config) error {
	l, err := InitWithConfig(ctx, config, WithWriter(os.Stdout))
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
