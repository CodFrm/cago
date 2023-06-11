package producer

import (
	"context"
	"encoding/json"

	"github.com/codfrm/cago/pkg/broker"
	broker2 "github.com/codfrm/cago/pkg/broker/broker"
)

type ExampleMsg struct {
	Time int64
}

// PublishExample 发布示例消息
func PublishExample(ctx context.Context, msg *ExampleMsg) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return broker.Default().Publish(ctx, ExampleTopic, &broker2.Message{
		Body: body,
	})
}

func SubscribeExample(ctx context.Context, fn func(ctx context.Context, msg *ExampleMsg) error) error {
	_, err := broker.Default().Subscribe(ctx, ExampleTopic, func(ctx context.Context, event broker2.Event) error {
		msg, err := ParseExampleMsg(event.Message())
		if err != nil {
			return err
		}
		return fn(ctx, msg)
	}, broker2.Retry())
	return err
}

func ParseExampleMsg(msg *broker2.Message) (*ExampleMsg, error) {
	ret := &ExampleMsg{}
	if err := json.Unmarshal(msg.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}
