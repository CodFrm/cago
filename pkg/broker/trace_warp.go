package broker

import (
	"context"
	"encoding/json"

	broker2 "github.com/codfrm/cago/pkg/broker/broker"
	"github.com/codfrm/cago/pkg/logger"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type traceBroker struct {
	tracer trace.Tracer
	wrap   broker2.Broker
}

// wrapTrace 包装链路追踪
func wrapTrace(broker broker2.Broker, tracer trace.Tracer) broker2.Broker {
	return &traceBroker{wrap: broker, tracer: tracer}
}

func (t *traceBroker) Publish(ctx context.Context, topic string, data *broker2.Message, opts ...broker2.PublishOption) error {
	ctx, span := t.tracer.Start(ctx, "send "+topic,
		trace.WithAttributes(
			attribute.String("messaging.system", t.wrap.String()),
			attribute.String("messaging.destination", topic),
			attribute.String("messaging.destination_kind", "queue"),
		),
		trace.WithSpanKind(trace.SpanKindProducer),
	)
	defer span.End()
	bt, err := span.SpanContext().MarshalJSON()
	if err != nil {
		return err
	}
	if data.Header == nil {
		data.Header = make(map[string]string)
	}
	data.Header["spanConfig"] = string(bt)
	return t.wrap.Publish(ctx, topic, data, opts...)
}

type spanContextConfig struct {
}

func (t *traceBroker) Subscribe(ctx context.Context, topic string, h broker2.Handler, opts ...broker2.SubscribeOption) (broker2.Subscriber, error) {
	options := broker2.NewSubscribeOptions(opts...)
	return t.wrap.Subscribe(ctx, topic, func(ctx context.Context, event broker2.Event) error {
		spanConfig := trace.SpanContextConfig{}
		if s, ok := event.Message().Header["spanConfig"]; ok {
			if err := json.Unmarshal([]byte(s), &spanConfig); err == nil {
				spanCtx := trace.NewSpanContext(spanConfig)
				ctx = trace.ContextWithRemoteSpanContext(ctx, spanCtx)
			} else {
				logger.Ctx(ctx).Error("broker subscribe unmarshal spanConfig error", zap.Error(err))
			}
		}
		ctx, span := t.tracer.Start(ctx, "process "+topic,
			trace.WithAttributes(
				attribute.String("messaging.system", t.wrap.String()),
				attribute.String("messaging.operation", "process"),
				attribute.String("messaging.destination", topic),
				attribute.String("messaging.destination_kind", "queue"),
				attribute.String("messaging.group", options.Group),
			),
			trace.WithSpanKind(trace.SpanKindConsumer),
		)
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
