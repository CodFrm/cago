package mux

import (
	"net/http"

	"github.com/codfrm/cago/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type customResponseWriter struct {
	gin.ResponseWriter
	status int
}

func (c *customResponseWriter) WriteHeader(status int) {
	c.status = status
	c.ResponseWriter.WriteHeader(status)
}

func initHandler(c *gin.Context) {
	logger := logger.Default().With(
		zap.String("request_id", uuid.New().String()),
		zap.String("client_ip", c.ClientIP()),
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("user_agent", c.Request.UserAgent()),
	)
	c.Set("logger", logger)
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				logger.Error("internal server error", zap.Error(err), zap.Stack("stack"))
				_ = c.AbortWithError(http.StatusInternalServerError, err)
			} else {
				panic(err)
			}
		}
	}()

	custom := &customResponseWriter{ResponseWriter: c.Writer}
	c.Writer = custom
	// 处理错误日志
	c.Next()
}
