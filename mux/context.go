package mux

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type WebContext struct {
	*gin.Context
	logger *zap.Logger
}

func NewContext(ctx *gin.Context, logger *zap.Logger) *WebContext {
	return &WebContext{Context: ctx, logger: logger}
}

func (c *WebContext) Error(msg string, fields ...zap.Field) {
	c.logger.Error(msg, fields...)
}

func (c *WebContext) Logger() *zap.Logger {
	return c.logger
}
