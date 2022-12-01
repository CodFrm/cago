package broker

import (
	"context"

	"github.com/codfrm/cago/configs"
	broker2 "github.com/codfrm/cago/pkg/broker/broker"
	"github.com/codfrm/cago/pkg/trace"
)

var broker broker2.Broker

func Broker(ctx context.Context, config *configs.Config) error {
	cfg := &Config{}
	if err := config.Scan("broker", cfg); err != nil {
		return err
	}
	cfg.defaultGroup = config.AppName
	cfg.topicPrefix = string(config.Env)
	options := make([]Option, 0)
	if tp := trace.Default(); tp != nil {
		options = append(options, WithTracer(tp.Tracer(config.AppName+".broker")))
	}
	b, err := NewWithConfig(ctx, cfg, options...)
	if err != nil {
		return err
	}
	broker = b
	return nil
}

func SetBroker(b broker2.Broker) {
	broker = b
}

func Default() broker2.Broker {
	return broker
}
