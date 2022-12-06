package broker

import (
	"context"

	"github.com/codfrm/cago"
	"github.com/codfrm/cago/configs"
	broker2 "github.com/codfrm/cago/pkg/broker/broker"
	"github.com/codfrm/cago/pkg/trace"
	trace2 "go.opentelemetry.io/otel/trace"
)

var broker broker2.Broker

const instrumName = "github.com/codfrm/cago/pkg/broker"

func Broker(ctx context.Context, config *configs.Config) error {
	cfg := &Config{}
	if err := config.Scan("broker", cfg); err != nil {
		return err
	}
	cfg.defaultGroup = config.AppName
	cfg.topicPrefix = string(config.Env)
	options := make([]Option, 0)
	if tp := trace.Default(); tp != nil {
		options = append(options, WithTracer(tp.Tracer(
			instrumName,
			trace2.WithInstrumentationVersion("semver:"+cago.Version()),
		)))
	}
	b, err := NewWithConfig(ctx, cfg, options...)
	if err != nil {
		return err
	}
	broker = b
	return nil
}

func WithCallback(callback func(ctx context.Context, broker broker2.Broker) error) cago.FuncComponent {
	return func(ctx context.Context, cfg *configs.Config) error {
		err := Broker(ctx, cfg)
		if err != nil {
			return err
		}
		return callback(ctx, broker)
	}
}

func SetBroker(b broker2.Broker) {
	broker = b
}

func Default() broker2.Broker {
	return broker
}
