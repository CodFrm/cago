package token

import (
	"github.com/codfrm/cago/pkg/utils"
	"github.com/gin-gonic/gin"
)

func UserId(ctx *gin.Context) (int64, bool) {
	u, ok := ctx.Get(Userinfo)
	if !ok {
		return 0, false
	}
	return utils.ToNumber[int64](u.(gin.H)["uid"].(string)), true
}

func Isadmin(ctx *gin.Context) (int64, bool) {
	u, ok := ctx.Get(Userinfo)
	if !ok {
		return 0, false
	}
	return utils.ToNumber[int64](u.(gin.H)["uid"].(string)), false
}

func Authtoken(ctx *gin.Context) (*Token, bool) {
	t, ok := ctx.Get(AuthToken)
	if !ok {
		return nil, false
	}
	return t.(*Token), true
}
