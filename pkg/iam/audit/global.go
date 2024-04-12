package audit

import (
	"context"
	"go.uber.org/zap"
)

type key int

const (
	fieldsKey key = iota
)

var defaultAudit *Audit

func SetDefault(audit *Audit) {
	defaultAudit = audit
}

func Default() *Audit {
	return defaultAudit
}

func WithFields(ctx context.Context, fields ...zap.Field) context.Context {
	ctxFields, ok := ctx.Value(fieldsKey).([]zap.Field)
	if ok {
		ctxFields = append(ctxFields, fields...)
	} else {
		ctxFields = fields
	}
	return context.WithValue(ctx, fieldsKey, ctxFields)
}
