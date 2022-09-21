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
}

func newSubscribe(b *nsqBroker, topic string, handler broker.Handler, options broker.SubscribeOptions) (broker.Subscriber, error) {
	consumer, err := nsq.NewConsumer(topic, options.Group, b.config)
	if err != nil {
		return nil, err
	}
	ret := &subscribe{
		consumer: consumer, handler: handler,
	}
	logger := logger.Default().With(zap.String("topic", topic))
	ret.consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) (err error) {
		data := &broker.Message{}
		defer func() {
			if err == nil {
				if options.AutoAck {
					message.Finish()
				}
			} else {
				message.Requeue(-1)
				logger.Error("nsq subscriber handle error", zap.Error(err))
			}
		}()
		if err = json.Unmarshal(message.Body, data); err != nil {
			logger.Error("nsq subscriber unmarshal error", zap.Error(err))
			return err
		}
		err = handler(context.Background(), &event{
			topic:   topic,
			data:    data,
			message: message,
		})
		return err
	}))
	if err := ret.consumer.ConnectToNSQD(b.address); err != nil {
		return nil, err
	}
	return ret, nil
}

func (s *subscribe) Topic() string {
	return ""
}

func (s *subscribe) Unsubscribe() error {
	return nil
}
