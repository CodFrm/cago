package logger

import (
	"context"

	"go.uber.org/zap"
)

func Ctx(ctx context.Context) *zap.Logger {
	return Logger
}
