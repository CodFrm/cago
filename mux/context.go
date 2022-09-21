package mux

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Context struct {
	*gin.Context
	logger *zap.Logger
}

func NewContext(ctx *gin.Context, logger *zap.Logger) *Context {
	return &Context{Context: ctx, logger: logger}
}

func (c *Context) Logger() *zap.Logger {
	return c.logger
}
