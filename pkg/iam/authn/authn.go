package authn

import (
	"context"

	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/gin-gonic/gin"
)

type key int

const (
	iamSession key = iota
)

type Authn struct {
	options  *Options
	database Database
}

func New(database Database, opts ...Option) *Authn {
	options := newOptions(opts...)

	return &Authn{
		options:  options,
		database: database,
	}
}

// Middleware 中间件
// force 表示是否强制要求认证
func (a *Authn) Middleware(force bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session, err := a.options.sessionManager.GetFromRequest(ctx)
		if err != nil {
			httputils.HandleResp(ctx, err)
			return
		}
		if session == nil {
			if force {
				httputils.HandleResp(ctx, ErrUnauthorized)
				return
			}
			return
		}
		reqCtx := ctx.Request.Context()
		if a.options.setContext != nil {
			reqCtx, err = a.options.setContext(reqCtx, session)
			if err != nil {
				httputils.HandleResp(ctx, err)
				return
			}
		}
		ctx.Request = ctx.Request.WithContext(context.WithValue(reqCtx, iamSession, session))
	}
}
