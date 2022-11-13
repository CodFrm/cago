package httputils

import (
	"net/http"

	"github.com/codfrm/cago/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	pkgValidator "github.com/scriptscat/cloudcat/pkg/utils/validator"
	"go.uber.org/zap"
)

type List struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
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

func HandleResp(ctx *gin.Context, resp interface{}) {
	switch resp.(type) {
	case *JsonResponseError:
		err := resp.(*JsonResponseError)
		ctx.JSON(err.Status, err)
	case validator.ValidationErrors:
		err := resp.(validator.ValidationErrors)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": -1, "msg": pkgValidator.TransError(err),
		})
	case error:
		err := resp.(error)
		logger := logger.Ctx(ctx).With(
			zap.String("url", ctx.Request.URL.String()),
			zap.String("method", ctx.Request.Method),
			zap.String("ip", ctx.ClientIP()),
		)
		logger.Error("internal server error", zap.Error(err), zap.StackSkip("stack", 3))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": -1, "msg": "系统错误",
		})
	case *List:
		list := resp.(*List)
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0, "msg": "success", "list": list.List, "total": list.Total,
		})
	default:
		ctx.JSON(http.StatusOK, JsonResponse{
			Code: 0,
			Msg:  "success",
			Data: resp,
		})
	}
}
