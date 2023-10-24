package utils

import (
	"context"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/pkg/opentelemetry/trace"
)

// BaseContext 组装基础context, 例如logger, trace等
func BaseContext(parentCtx context.Context) context.Context {
	ctx := context.Background()
	ctx = logger.ContextWithLogger(ctx, logger.Ctx(parentCtx))
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		ctx = trace.ContextWithSpan(ctx, span)
	}
	return ctx
}
