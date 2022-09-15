package broker

import (
	"context"

	"github.com/codfrm/cago/configs"
	broker2 "github.com/codfrm/cago/pkg/broker/broker"
	"github.com/codfrm/cago/pkg/trace"
)

var broker broker2.Broker

func Broker(ctx context.Context, config *configs.Config) error {
	b, err := InitWithConfig(ctx, config, WithTraceProvider(trace.Default()))
	if err != nil {
		return err
	}
	broker = b
	return nil
}

func Default() broker2.Broker {
	return broker
}
