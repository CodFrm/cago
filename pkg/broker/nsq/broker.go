package nsq

import (
	"context"
	"encoding/json"

	"github.com/codfrm/cago/pkg/broker/broker"
	"github.com/nsqio/go-nsq"
)

type nsqBroker struct {
	address  string
	config   *nsq.Config
	producer *nsq.Producer
}

func NewBroker(address string) (broker.Broker, error) {
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer(address, config)
	if err != nil {
		return nil, err
	}
	return &nsqBroker{
		address:  address,
		config:   config,
		producer: producer,
	}, nil
}

func (b *nsqBroker) Publish(ctx context.Context, topic string, data *broker.Message, opt ...broker.PublishOption) error {
	bt, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return b.producer.Publish(topic, bt)
}

func (b *nsqBroker) Subscribe(ctx context.Context, topic string, h broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	options := broker.NewSubscribeOptions(opts...)
	return newSubscribe(b, topic, h, options)
}

func (b *nsqBroker) Close() error {
	b.producer.Stop()
	return nil
}

func (b *nsqBroker) String() string {
	return "nsq"
}
