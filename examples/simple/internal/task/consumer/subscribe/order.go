package subscribe

import (
	"context"
	"github.com/codfrm/cago/examples/simple/internal/task/producer"
	"github.com/codfrm/cago/pkg/logger"
	"go.uber.org/zap"
)

type Example struct {
}

func (u *Example) Subscribe(ctx context.Context) error {
	if err := producer.SubscribeExample(ctx, u.example); err != nil {
		return err
	}
	return nil
}

func (u *Example) example(ctx context.Context, msg *producer.ExampleMsg) error {
	logger.Ctx(ctx).Info("收到消息", zap.Int64("time", msg.Time))
	return nil
}
