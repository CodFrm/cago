package audit

import (
	"context"
)

type key int

const (
	auditKey key = iota
)

var defaultAudit *Audit

func SetDefault(audit *Audit) {
	defaultAudit = audit
}

func Default() *Audit {
	return defaultAudit
}

func Ctx(ctx context.Context) *CtxAudit {
	audit, ok := ctx.Value(auditKey).(*Audit)
	if ok {
		return NewCtxAudit(ctx, audit)
	}
	return NewCtxAudit(ctx, defaultAudit)
}

func WithAudit(ctx context.Context, audit *Audit) context.Context {
	return context.WithValue(ctx, auditKey, audit)
}
