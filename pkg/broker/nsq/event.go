package nsq

import (
	"time"

	"github.com/codfrm/cago/pkg/broker/broker"
	"github.com/nsqio/go-nsq"
)

type event struct {
	topic     string
	data      *broker.Message
	message   *nsq.Message
	isRequeue bool
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

func (e *event) Requeue(delay time.Duration) error {
	e.message.Requeue(delay)
	e.isRequeue = true
	return nil
}
