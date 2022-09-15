package broker

import (
	"context"
	"encoding/json"

	broker2 "github.com/codfrm/cago/pkg/broker/broker"
	"go.opentelemetry.io/otel/trace"
)

type traceBroker struct {
	tracer trace.Tracer
	wrap   broker2.Broker
}

// wrapTrace 包装链路追踪
func wrapTrace(name string, broker broker2.Broker, tracerProvider trace.TracerProvider) broker2.Broker {
	return &traceBroker{wrap: broker, tracer: tracerProvider.Tracer(name)}
}

func (t *traceBroker) Publish(ctx context.Context, topic string, data *broker2.Message, opts ...broker2.PublishOption) error {
	ctx, span := t.tracer.Start(ctx, "broker.Publish")
	defer span.End()
	bt, err := span.SpanContext().MarshalJSON()
	if err != nil {
		return err
	}
	data.Header["spanConfig"] = string(bt)
	return t.wrap.Publish(ctx, topic, data, opts...)
}

func (t *traceBroker) Subscribe(ctx context.Context, topic string, h broker2.Handler, opts ...broker2.SubscribeOption) (broker2.Subscriber, error) {
	return t.wrap.Subscribe(ctx, topic, func(ctx context.Context, event broker2.Event) error {
		spanConfig := trace.SpanContextConfig{}
		if s, ok := event.Message().Header["spanConfig"]; ok {
			if err := json.Unmarshal([]byte(s), &spanConfig); err == nil {
				spanCtx := trace.NewSpanContext(spanConfig)
				ctx = trace.ContextWithRemoteSpanContext(ctx, spanCtx)
			}
		}
		ctx, span := t.tracer.Start(ctx, "broker.Subscribe")
		defer span.End()
		return h(ctx, event)
	}, opts...)
}

func (t *traceBroker) Close() error {
	return t.wrap.Close()
}

func (t *traceBroker) String() string {
	return t.wrap.String()
}
