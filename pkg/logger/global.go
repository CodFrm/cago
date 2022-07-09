package logger

import (
	"context"

	"github.com/codfrm/cago"
	"github.com/codfrm/cago/config"
	"go.uber.org/zap"
)

var logger *zap.Logger

// Logger 日志组件,默认是会注册到cago的
func Logger(ctx context.Context, config *config.Config) error {
	l, err := InitWithConfig(ctx, config)
	if err != nil {
		return err
	}
	logger = l.With(zap.String("app", config.AppName))
	return nil
}

func SetLogger(l *zap.Logger) {
	logger = l
}

func Ctx(ctx cago.Context) *zap.Logger {
	return logger
}
