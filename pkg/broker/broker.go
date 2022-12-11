package broker

import (
	"context"
	"errors"

	broker2 "github.com/codfrm/cago/pkg/broker/broker"
	"github.com/codfrm/cago/pkg/broker/event_bus"
	"github.com/codfrm/cago/pkg/broker/nsq"
)

type Type string

const (
	NSQ      Type = "nsq"
	EventBus Type = "event_bus"
)

type Config struct {
	Type Type
	NSQ  nsq.Config
}

func NewWithConfig(ctx context.Context, cfg *Config, opts ...Option) (broker2.Broker, error) {
	var ret broker2.Broker
	var err error
	switch cfg.Type {
	case NSQ:
		ret, err = nsq.NewBroker(cfg.NSQ)
	case EventBus:
		ret = event_bus.NewEvBusBroker()
	default:
		return nil, errors.New("type not found")
	}
	if err != nil {
		return nil, err
	}
	opts = append(opts, WithBroker(ret))
	return New(opts...)
}

func New(opts ...Option) (broker2.Broker, error) {
	options := &Options{}
	for _, o := range opts {
		o(options)
	}
	ret := options.broker
	if options.tracer != nil {
		ret = wrapTrace(ret, options.tracer)
	}
	return newWrap(ret, options), nil
}
