package event_bus

import (
	"github.com/codfrm/cago/pkg/broker/broker"
)

type event struct {
	topic string
	data  *broker.Message
}

func (e *event) Topic() string {
	return e.topic
}

func (e *event) Message() *broker.Message {
	return e.data
}

func (e *event) Ack() error {
	return nil
}

func (e *event) Error() error {
	return nil
}
