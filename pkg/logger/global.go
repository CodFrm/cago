package logger

import (
	"context"

	"github.com/codfrm/cago"
	"github.com/codfrm/cago/configs"
	"go.uber.org/zap"
)

var logger *zap.Logger

// Logger 日志组件,核心组件,必须注册
func Logger(ctx context.Context, config *configs.Config) error {
	l, err := InitWithConfig(ctx, config, AppendLabels(zap.String("app", config.AppName)))
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

func Default() *zap.Logger {
	return logger
}
