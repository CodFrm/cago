package broker

import (
	"context"
	"errors"

	broker2 "github.com/codfrm/cago/pkg/broker/broker"
	"github.com/codfrm/cago/pkg/broker/event_bus"
	"github.com/codfrm/cago/pkg/broker/nsq"
)

type BrokerType string

const (
	NSQ       BrokerType = "nsq"
	EVENT_BUS BrokerType = "event_bus"
)

type NSQConfig struct {
	Addr string
}

type Config struct {
	Type BrokerType
	NSQ  *NSQConfig
}

func InitWithConfig(ctx context.Context, cfg *Config, opts ...Option) (broker2.Broker, error) {
	var ret broker2.Broker
	var err error
	switch cfg.Type {
	case NSQ:
		ret, err = nsq.NewBroker(cfg.NSQ.Addr)
	case EVENT_BUS:
		ret = event_bus.NewEvBusBroker()
	default:
		return nil, errors.New("type not found")
	}
	if err != nil {
		return nil, err
	}
	opts = append(opts, WithBroker(ret))
	return Init(opts...)
}

func Init(opts ...Option) (broker2.Broker, error) {
	options := &Options{}
	for _, o := range opts {
		o(options)
	}
	ret := options.broker
	if options.tracer != nil {
		ret = wrapTrace(ret, options.tracer)
	}
	return ret, nil
}
