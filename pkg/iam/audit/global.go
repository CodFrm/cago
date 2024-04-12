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

func Ctx(ctx context.Context) *Audit {
	audit, ok := ctx.Value(auditKey).(*Audit)
	if ok {
		return audit
	}
	return defaultAudit
}

func WithAudit(ctx context.Context, audit *Audit) context.Context {
	return context.WithValue(ctx, auditKey, audit)
}
