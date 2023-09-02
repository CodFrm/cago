package httputils

import (
	"github.com/codfrm/cago/pkg/opentelemetry/trace"
	"net/http"

	"github.com/codfrm/cago/pkg/errs"

	"github.com/codfrm/cago/pkg/logger"
	pkgValidator "github.com/codfrm/cago/pkg/utils/validator"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

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
		ctx.AbortWithStatusJSON(data.Status, data)
	case validator.ValidationErrors:
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": -1, "msg": pkgValidator.TransError(data),
		})
	case error:
		requestId := trace.RequestID(ctx)
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
