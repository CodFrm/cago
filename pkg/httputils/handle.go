package httputils

import (
	"net/http"

	"github.com/codfrm/cago/mux"
	"github.com/codfrm/cago/pkg/errs"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	pkgValidator "github.com/scriptscat/cloudcat/pkg/utils/validator"
	"go.uber.org/zap"
)

type List struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

func Handle(ctx *mux.WebContext, f func() interface{}) {
	resp := f()
	if resp == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0, "msg": "ok",
		})
		return
	}
	switch resp.(type) {
	case *errs.JsonRespondError:
		err := resp.(*errs.JsonRespondError)
		ctx.JSON(err.Status, err)
	case validator.ValidationErrors:
		err := resp.(validator.ValidationErrors)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": -1, "msg": pkgValidator.TransError(err),
		})
	case error:
		err := resp.(error)
		ctx.Error("server internal error", zap.Error(err), zap.Stack("stack"))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": -1, "msg": "系统错误",
		})
	case *List:
		list := resp.(*List)
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0, "msg": "ok", "list": list.List, "total": list.Total,
		})
	default:
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0, "msg": "ok", "data": resp,
		})
	}
}

func HandleError(ctx *mux.WebContext, err error) {
	switch err.(type) {
	case *errs.JsonRespondError:
		err := err.(*errs.JsonRespondError)
		ctx.JSON(err.Status, err)
	case validator.ValidationErrors:
		err := err.(validator.ValidationErrors)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": -1, "msg": pkgValidator.TransError(err),
		})
	case error:
		err := err.(error)
		ctx.Error("server internal error", zap.Error(err), zap.Stack("stack"))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": -1, "msg": "系统错误",
		})
	}
}
