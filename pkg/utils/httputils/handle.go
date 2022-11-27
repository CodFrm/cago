package httputils

import (
	"net/http"

	"github.com/codfrm/cago/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	pkgValidator "github.com/scriptscat/cloudcat/pkg/utils/validator"
	"go.uber.org/zap"
)

type Page struct {
	Page  int `form:"page" binding:"required"`
	Limit int `form:"limit" binding:"required"`
}

func (p *Page) GetPage() int {
	if p.Page == 0 {
		return 1
	}
	return p.Page
}

func (p *Page) GetOffset() int {
	return (p.GetPage() - 1) * p.Limit
}

func (p *Page) GetLimit() int {
	if p.Limit == 0 {
		return 20
	}
	return p.Limit
}

type List[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
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
	case *List:
		list := resp
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
