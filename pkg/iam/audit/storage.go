package audit

import (
	"context"

	"go.uber.org/zap"
)

type Storage interface {
	Record(ctx context.Context, module, eventName string, fields ...zap.Field) error
}
