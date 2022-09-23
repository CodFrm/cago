package logger

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type loggerContextKeyType int

const loggerKey loggerContextKeyType = iota

func Middleware(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger := logger.With(
			zap.String("client_ip", ctx.ClientIP()),
			zap.String("method", ctx.Request.Method),
			zap.String("path", ctx.Request.URL.Path),
			zap.String("user_agent", ctx.Request.UserAgent()),
		)
		ctx.Request = ctx.Request.WithContext(ContextWithLogger(ctx.Request.Context(), logger))
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					Ctx(ctx).Error("internal server error", zap.Error(err), zap.Stack("stack"))
					ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"code": -1000,
						"msg":  "internal server error",
					})
				} else {
					panic(err)
				}
			}
		}()
		// 处理错误日志
		ctx.Next()
	}
}

func Ctx(ctx context.Context) *zap.Logger {
	if gctx, ok := ctx.(*gin.Context); ok {
		return gctx.Request.Context().Value(loggerKey).(*zap.Logger)
	}
	return ctx.Value(loggerKey).(*zap.Logger)
}

func ContextWithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}
