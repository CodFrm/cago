package broker

import (
	"context"

	"github.com/codfrm/cago"
	"github.com/codfrm/cago/configs"
	broker2 "github.com/codfrm/cago/pkg/broker/broker"
	"github.com/codfrm/cago/pkg/opentelemetry/trace"
	trace2 "go.opentelemetry.io/otel/trace"
)

var broker broker2.Broker

const instrumName = "github.com/codfrm/cago/pkg/broker"

func Broker(ctx context.Context, config *configs.Config) error {
	cfg := &Config{}
	if err := config.Scan("broker", cfg); err != nil {
		return err
	}
	options := make([]Option, 0)
	if tp := trace.Default(); tp != nil {
		options = append(options, WithTracer(tp.Tracer(
			instrumName,
			trace2.WithInstrumentationVersion("semver:"+cago.Version()),
		)))
	}
	options = append(options, WithDefaultGroup(config.AppName), WithTopicPrefix(string(config.Env)))
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
