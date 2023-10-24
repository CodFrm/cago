package gogo

import (
	"context"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/pkg/opentelemetry/trace"
	"sync"
)

var wg sync.WaitGroup

// Go 框架处理协程,用于优雅启停
func Go(fun func(ctx context.Context) error, opts ...Option) error {
	wg.Add(1)
	options := &Options{}
	for _, o := range opts {
		o(options)
	}
	if options.ctx == nil {
		options.ctx = context.Background()
	}
	go func() {
		defer wg.Done()
		_ = fun(options.ctx)
	}()
	return nil
}

// Wait 等待所有协程结束
func Wait() {
	wg.Wait()
}

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
