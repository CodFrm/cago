package api

import (
	"context"
	_ "github.com/codfrm/cago/examples/simple/docs"
	"github.com/codfrm/cago/examples/simple/internal/controller/example_ctr"
	"github.com/codfrm/cago/examples/simple/internal/controller/user_ctr"
	"github.com/codfrm/cago/pkg/iam/audit"
	"github.com/codfrm/cago/pkg/iam/authn"
	"github.com/codfrm/cago/server/mux"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Router 路由
// @title    api文档
// @version  1.0
// @BasePath /api/v1
func Router(ctx context.Context, root *mux.Router) error {
	r := root.Group("/api/v1")

	userLoginCtr := user_ctr.NewUser()
	{
		// 绑定路由
		r.Group("/").Bind(
			userLoginCtr.Register,
			userLoginCtr.Login,
		)

		r.Group("/", authn.Default().Middleware(true)).Bind(
			userLoginCtr.CurrentUser,
			userLoginCtr.Logout,
			userLoginCtr.RefreshToken,
		)
	}

	{
		exampleCtl := example_ctr.NewExample()
		r.Group("/").Bind(
			exampleCtl.Ping,
			exampleCtl.GinFun,
		)

		r.Group("/", authn.Default().Middleware(true), audit.Default().Middleware(func(ctx *gin.Context) []zap.Field {
			return []zap.Field{
				zap.String("path", ctx.Request.URL.Path),
			}
		})).Bind(
			exampleCtl.Audit,
		)
	}

	return nil
}
