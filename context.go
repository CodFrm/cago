package cago

import (
	"context"

	"go.uber.org/zap"
)

type Context interface {
	context.Context
	Logger() *zap.Logger
}

type backgroundContext struct {
	context.Context
	logger *zap.Logger
}

func (b *backgroundContext) Logger() *zap.Logger {
	return nil
}

func Background() Context {
	return &backgroundContext{Context: context.Background()}
}
