package httputils

import (
	"net/http"

	"github.com/codfrm/cago/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	pkgValidator "github.com/scriptscat/cloudcat/pkg/utils/validator"
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

func HandleResp(ctx *gin.Context, resp interface{}) {
	switch resp := resp.(type) {
	case *JsonResponseError:
		err := resp
		ctx.AbortWithStatusJSON(err.Status, err)
	case validator.ValidationErrors:
		err := resp
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": -1, "msg": pkgValidator.TransError(err),
		})
	case error:
		err := resp
		logger := logger.Ctx(ctx).With(
			zap.String("url", ctx.Request.URL.String()),
			zap.String("method", ctx.Request.Method),
			zap.String("ip", ctx.ClientIP()),
		)
		logger.Error("internal server error", zap.Error(err), zap.StackSkip("stack", 3))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code": -1, "msg": "系统错误",
		})
	default:
		ctx.JSON(http.StatusOK, JsonResponse{
			Code: 0,
			Msg:  "success",
			Data: resp,
		})
	}
}
