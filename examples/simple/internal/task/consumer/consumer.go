package consumer

import (
	"context"
	"github.com/codfrm/cago"
	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/examples/simple/internal/task/consumer/subscribe"
)

type Subscribe interface {
	Subscribe(ctx context.Context) error
}

// Consumer 消费者
func Consumer() cago.FuncComponent {
	return func(ctx context.Context, cfg *configs.Config) error {
		subscribers := []Subscribe{
			&subscribe.Example{},
		}
		for _, v := range subscribers {
			if err := v.Subscribe(ctx); err != nil {
				return err
			}
		}
		return nil
	}
}
