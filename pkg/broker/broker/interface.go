package broker

import (
	"context"
)

type Message struct {
	Header map[string]string
	Body   []byte
}

type Event interface {
	Topic() string
	Message() *Message
	Ack() error
	Error() error
}

type Handler func(ctx context.Context, event Event) error

type Broker interface {
	Publish(ctx context.Context, topic string, data *Message, opts ...PublishOption) error
	Subscribe(ctx context.Context, topic string, h Handler, opts ...SubscribeOption) (Subscriber, error)
	Close() error
	String() string
}

type Subscriber interface {
	Topic() string
	Unsubscribe() error
}
