package nsq

import (
	"context"
	"encoding/json"

	"github.com/codfrm/cago/pkg/broker/broker"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"
)

type subscribe struct {
	consumer *nsq.Consumer
	handler  broker.Handler
	topic    string
	config   *Config
}

func newSubscribe(b *nsqBroker, topic string, handler broker.Handler, options broker.SubscribeOptions) (broker.Subscriber, error) {
	consumer, err := nsq.NewConsumer(topic, options.Group, b.nsqConfig)
	if err != nil {
		return nil, err
	}
	ret := &subscribe{
		consumer: consumer, handler: handler,
		topic: topic, config: b.config,
	}
	logger := logger.Default().With(zap.String("topic", topic), zap.String("group", options.Group))
	ret.consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) (err error) {
		message.DisableAutoResponse()
		data := &broker.Message{}
		ev := &event{
			topic:     topic,
			data:      data,
			message:   message,
			attempted: int(message.Attempts),
		}
		defer func() {
			if err != nil {
				if !ev.isAct {
					if options.Retry {
						message.Requeue(-1)
						logger.Error("nsq subscriber handle error", zap.Bool("retry", true), zap.Error(err))
					} else if options.AutoAck {
						message.Finish()
					}
				} else {
					logger.Error("nsq subscriber handle error", zap.Bool("retry", true), zap.Error(err))
				}
				logger.Error("nsq subscriber handle error", zap.Error(err))
			} else if options.AutoAck && !ev.isAct {
				message.Finish()
			}
		}()
		if err = json.Unmarshal(message.Body, data); err != nil {
			logger.Error("nsq subscriber unmarshal error", zap.Error(err))
			return err
		}
		err = handler(context.Background(), ev)
		return err
	}))
	if b.config.NSQLookupAddr != nil {
		err = ret.consumer.ConnectToNSQLookupds(b.config.NSQLookupAddr)
	} else {
		err = ret.consumer.ConnectToNSQD(b.config.Addr)
	}
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (s *subscribe) Topic() string {
	return s.topic
}

func (s *subscribe) Unsubscribe() error {
	if s.config.NSQLookupAddr != nil {
		for _, addr := range s.config.NSQLookupAddr {
			if err := s.consumer.DisconnectFromNSQLookupd(addr); err != nil {
				return err
			}
		}
		return nil
	}
	return s.consumer.DisconnectFromNSQD(s.config.Addr)
}
