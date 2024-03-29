package httputils

import (
	"context"
	"net/http"

	"github.com/codfrm/cago/pkg/errs"
	"github.com/codfrm/cago/pkg/logger"
	pkgValidator "github.com/codfrm/cago/pkg/utils/validator"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Handle 处理请求
func Handle(ctx *gin.Context, f func() interface{}) {
	resp := f()
	if resp == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0, "msg": "ok",
		})
		return
	}
	HandleResp(ctx, resp)
}

func deal(ctx *gin.Context, resp any, field []zap.Field) {
	// 从trace中获取
	switch data := resp.(type) {
	case *Error:
		data.RequestID = RequestID(ctx)
		ctx.AbortWithStatusJSON(data.Status, data)
	case validator.ValidationErrors:
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": -1, "msg": pkgValidator.TransError(data),
		})
	case error:
		requestId := RequestID(ctx)
		field = append(field, zap.Error(data))
		logger.Ctx(ctx).Error(
			"internal server error",
			field...,
		)
		if requestId != "" {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code": -1, "msg": "系统错误", "request_id": requestId,
			})
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code": -1, "msg": "系统错误",
			})
		}
	default:
		ctx.JSON(http.StatusOK, JSONResponse{
			Code: 0,
			Msg:  "success",
			Data: data,
		})
	}
}

// HandleResp 处理响应
// 1. 如果是 httputils.Error 类型，会根据里面定义的状态码和数据返回
// 2. 如果是 validator.ValidationErrors 类型，会返回 400 错误码和错误信息
// 3. 如果是 error 类型，会返回 500 错误码，错误信息会记录到日志中
// 4. 其他情况会返回 200 状态码，并返回数据
func HandleResp(ctx *gin.Context, resp any) {
	var field []zap.Field
	for {
		if err, ok := resp.(errs.Unwrap); ok {
			switch err := err.(type) {
			case *errs.Error:
				field = append(field, err.Field()...)
			}
			resp = err.Unwrap()
		} else {
			break
		}
	}
	deal(ctx, resp, field)
}

// RequestID 获取请求ID
func RequestID(ctx context.Context) string {
	if span := trace.SpanFromContext(ctx); span != nil {
		if span.SpanContext().TraceID().IsValid() {
			return span.SpanContext().TraceID().String()
		}
	}
	return ""
}
