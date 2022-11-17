// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Based on https://github.com/DataDog/dd-trace-go/blob/8fb554ff7cf694267f9077ae35e27ce4689ed8b6/contrib/gin-gonic/gin/gintrace.go

package trace

import (
	"context"
	"fmt"

	"github.com/codfrm/cago/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type traceContextKeyType int

const tracerKey traceContextKeyType = iota

const (
	tracerName = "github.com/codfrm/cago/pkg/trace"
)

// Middleware 链路追踪中间件
func Middleware(serviceName string, tracerProvider trace.TracerProvider) gin.HandlerFunc {
	tracer := tracerProvider.Tracer(
		tracerName,
		trace.WithInstrumentationVersion("0.1.0"),
	)
	propagators := otel.GetTextMapPropagator()
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		ctx = propagators.Extract(ctx, propagation.HeaderCarrier(c.Request.Header))
		opts := []trace.SpanStartOption{
			trace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", c.Request)...),
			trace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(c.Request)...),
			trace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(serviceName, c.FullPath(), c.Request)...),
			trace.WithSpanKind(trace.SpanKindServer),
		}
		spanName := c.FullPath()
		if spanName == "" {
			spanName = fmt.Sprintf("HTTP %s route not found", c.Request.Method)
		}
		ctx, span := tracer.Start(ctx, spanName, opts...)
		defer span.End()

		// 给logger加上traceID
		ctx = logger.ContextWithLogger(ctx, logger.Ctx(c.Request.Context()).
			With(zap.String("trace_id", span.SpanContext().TraceID().String())))

		// 请求带上traceID
		c.Header("X-Trace-Id", span.SpanContext().TraceID().String())

		// 放入tracer
		ctx = context.WithValue(ctx, tracerKey, tracer)

		// pass the span through the request context
		c.Request = c.Request.WithContext(ctx)

		// serve the request to the next middleware
		c.Next()

		status := c.Writer.Status()
		attrs := semconv.HTTPAttributesFromHTTPStatusCode(status)
		spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCodeAndSpanKind(status, trace.SpanKindServer)
		span.SetAttributes(attrs...)
		span.SetStatus(spanStatus, spanMessage)
		if len(c.Errors) > 0 {
			span.SetAttributes(attribute.String("gin.errors", c.Errors.String()))
		}
	}
}

func SpanFromContext(ctx context.Context) trace.Span {
	if gctx, ok := ctx.(*gin.Context); ok {
		return trace.SpanFromContext(gctx.Request.Context())
	}
	return trace.SpanFromContext(ctx)
}

type warpTracer struct {
	trace.Tracer
}

func (w *warpTracer) Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	if gctx, ok := ctx.(*gin.Context); ok {
		return w.Tracer.Start(gctx.Request.Context(), spanName, opts...)
	}
	return w.Tracer.Start(ctx, spanName, opts...)
}

func TracerFromContext(ctx context.Context) trace.Tracer {
	if gctx, ok := ctx.(*gin.Context); ok {
		return &warpTracer{gctx.Request.Context().Value(tracerKey).(trace.Tracer)}
	}
	return ctx.Value(tracerKey).(trace.Tracer)
}