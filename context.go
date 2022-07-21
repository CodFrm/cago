package cago

import (
	"context"

	"go.uber.org/zap"
)

type Context interface {
	context.Context
	Logger() *zap.Logger
	UserID() int64
}

type backgroundContext struct {
	context.Context
}

func (b *backgroundContext) Logger() *zap.Logger {
	return nil
}

func (b *backgroundContext) UserID() int64 {
	return 0
}

func Background() Context {
	return &backgroundContext{Context: context.Background()}
}
