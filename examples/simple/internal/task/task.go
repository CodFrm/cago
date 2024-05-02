package task

import (
	"context"

	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/examples/simple/internal/task/crontab"
	"github.com/codfrm/cago/examples/simple/internal/task/queue/handler"
	"github.com/codfrm/cago/server/cron"
)

type Handler interface {
	Register(ctx context.Context) error
}

func Task(ctx context.Context, config *configs.Config) error {
	// 定时任务, 每5分钟执行一次
	_, err := cron.Default().AddFunc("*/5 * * * *", crontab.Example)
	if err != nil {
		return err
	}

	handlers := []Handler{
		&handler.Example{},
	}

	for _, h := range handlers {
		if err := h.Register(ctx); err != nil {
			return err
		}
	}

	return nil
}
