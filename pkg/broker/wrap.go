package broker

import (
	"context"

	broker2 "github.com/codfrm/cago/pkg/broker/broker"
)

type wrap struct {
	broker2.Broker
	options *Options
}

// newWrap 包装原有broker
func newWrap(broker broker2.Broker, options *Options) broker2.Broker {
	return &wrap{Broker: broker, options: options}
}

func (t *wrap) Publish(ctx context.Context, topic string, data *broker2.Message, opts ...broker2.PublishOption) error {
	if t.options.topicPrefix != "" {
		topic = t.options.topicPrefix + "." + topic
	}
	return t.Broker.Publish(ctx, topic, data, opts...)
}

func (t *wrap) Subscribe(ctx context.Context, topic string, h broker2.Handler, opts ...broker2.SubscribeOption) (broker2.Subscriber, error) {
	options := broker2.NewSubscribeOptions(opts...)
	if t.options.topicPrefix != "" {
		topic = t.options.topicPrefix + "." + topic
	}
	if t.options.defaultGroup != "" && options.Group == "" {
		opts = append(opts, broker2.Group(t.options.defaultGroup))
	}
	return t.Broker.Subscribe(ctx, topic, h, opts...)
}
