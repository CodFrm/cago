package nsq

import (
	"github.com/codfrm/cago/pkg/broker/broker"
	"github.com/nsqio/go-nsq"
)

type event struct {
	topic   string
	data    *broker.Message
	message *nsq.Message
}

func (e *event) Topic() string {
	return e.topic
}

func (e *event) Message() *broker.Message {
	return e.data
}

func (e *event) Ack() error {
	e.message.Finish()
	return nil
}

func (e *event) Error() error {
	return nil
}
