package broker

import (
	"context"
	"errors"

	"github.com/codfrm/cago/configs"
	broker2 "github.com/codfrm/cago/pkg/broker/broker"
	"github.com/codfrm/cago/pkg/broker/nsq"
)

type BrokerType string

const (
	NSQ BrokerType = "nsq"
)

type NSQConfig struct {
	Addr string
}

type Config struct {
	Type BrokerType
	NSQ  *NSQConfig
}

func InitWithConfig(ctx context.Context, config *configs.Config, opts ...Option) (broker2.Broker, error) {
	cfg := &Config{}
	if err := config.Scan("broker", cfg); err != nil {
		return nil, err
	}
	options := &Options{}
	for _, o := range opts {
		o(options)
	}
	var ret broker2.Broker
	var err error
	switch cfg.Type {
	case NSQ:
		ret, err = nsq.NewBroker(cfg.NSQ.Addr)
	default:
		return nil, errors.New("type not found")
	}
	if err != nil {
		return nil, err
	}
	if options.traceProvider != nil {
		ret = wrapTrace(config.AppName, ret, options.traceProvider)
	}
	return ret, nil
}
