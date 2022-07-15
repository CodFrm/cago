package logger

import (
	"context"

	"github.com/codfrm/cago"
	"github.com/codfrm/cago/config"
	"go.uber.org/zap"
)

var logger *zap.Logger

// Logger 日志组件,核心组件,必须注册
func Logger(ctx context.Context, config *config.Config) error {
	l, err := InitWithConfig(ctx, config, WithLabels(zap.String("app", config.AppName)))
	if err != nil {
		return err
	}
	logger = l
	return nil
}

func SetLogger(l *zap.Logger) {
	logger = l
}

func Ctx(ctx cago.Context) *zap.Logger {
	return logger
}
