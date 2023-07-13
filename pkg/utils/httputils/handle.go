package httputils

import (
	"github.com/codfrm/cago/pkg/utils/httputils/errs"
	"net/http"

	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/cago/pkg/logger"
	pkgValidator "github.com/codfrm/cago/pkg/utils/validator"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Unwrap interface {
	Unwrap() error
}

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
	switch data := resp.(type) {
	case *Error:
		ctx.AbortWithStatusJSON(data.Status, data)
	case validator.ValidationErrors:
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": -1, "msg": pkgValidator.TransError(data),
		})
	case *i18n.Error:
		ctx.AbortWithStatusJSON(data.Status(), gin.H{
			"code": data.Code(), "msg": data.Msg(i18n.DefaultLang),
		})
	case error:
		field = append(field, zap.Error(data))
		logger.Ctx(ctx).Error(
			"internal server error",
			field...,
		)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code": -1, "msg": "系统错误",
		})
	default:
		ctx.JSON(http.StatusOK, JsonResponse{
			Code: 0,
			Msg:  "success",
			Data: data,
		})
	}
}

func HandleResp(ctx *gin.Context, resp any) {
	var field []zap.Field
	for {
		if err, ok := resp.(Unwrap); ok {
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
