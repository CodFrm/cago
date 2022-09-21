package event_bus

import (
	"context"

	"github.com/codfrm/cago/pkg/broker/broker"
	"github.com/codfrm/cago/pkg/logger"
	"go.uber.org/zap"
)

type subscriber struct {
	e      *eventBusBroker
	topic  string
	handle func(data *broker.Message)
}

func newSubscriber(e *eventBusBroker, topic string, handler broker.Handler) (broker.Subscriber, error) {
	ret := &subscriber{
		e:     e,
		topic: topic,
	}
	logger := logger.Default().With(zap.String("topic", topic))
	ret.handle = func(data *broker.Message) {
		go func() {
			err := handler(context.Background(), &event{
				topic: topic,
				data:  data,
			})
			if err != nil {
				logger.Error("event bus subscriber handle error", zap.Error(err))
			}
		}()
	}
	if err := e.bus.Subscribe(topic, ret.handle); err != nil {
		return nil, err
	}
	return ret, nil
}

func (n *subscriber) Topic() string {
	return n.topic
}

func (n *subscriber) Unsubscribe() error {
	return n.e.bus.Unsubscribe(n.topic, n.handle)
}
