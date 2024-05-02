package crontab

import (
	"context"

	"github.com/codfrm/cago/pkg/logger"
)

func Example(ctx context.Context) error {
	logger.Ctx(ctx).Info("定时任务")
	return nil
}
