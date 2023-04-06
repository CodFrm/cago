package broker

import (
	"context"
	"time"
)

type Message struct {
	Header map[string]string
	Body   []byte
}

type Event interface {
	// Topic 主题
	Topic() string
	// Message 消息
	Message() *Message
	// Ack 确认
	Ack() error
	// Requeue 重新入队
	Requeue(delay time.Duration) error
	// Attempted 尝试次数
	Attempted() int
	// Error 错误
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
