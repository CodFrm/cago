package broker

import (
	"context"
	"time"
)

// Message 消息
type Message struct {
	// Header 消息头
	Header map[string]string
	// Body 消息体
	Body []byte
}

// Event 事件
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

// Broker 消息代理接口
type Broker interface {
	// Publish 发布一个消息，可以通过opts设置发布选项
	Publish(ctx context.Context, topic string, data *Message, opts ...PublishOption) error
	// Subscribe 订阅topic消息，可以通过opts设置订阅选项
	Subscribe(ctx context.Context, topic string, h Handler, opts ...SubscribeOption) (Subscriber, error)
	// Close 关闭消息代理
	Close() error
	// String 返回消息代理名称
	String() string
}

// Subscriber 订阅者接口
type Subscriber interface {
	// Topic 返回订阅的主题
	Topic() string
	// Unsubscribe 取消订阅
	Unsubscribe() error
}
