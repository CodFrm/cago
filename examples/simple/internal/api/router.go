package api

import (
	"context"

	_ "github.com/codfrm/cago/examples/simple/docs"
	"github.com/codfrm/cago/examples/simple/internal/controller/example_ctr"
	"github.com/codfrm/cago/examples/simple/internal/controller/user_ctr"
	"github.com/codfrm/cago/examples/simple/internal/repository/user_repo"
	"github.com/codfrm/cago/examples/simple/internal/service/user_svc"
	"github.com/codfrm/cago/pkg/iam/authn"
	"github.com/codfrm/cago/server/mux"
)

// Router 路由
// @title    api文档
// @version  1.0
// @BasePath /api/v1
func Router(ctx context.Context, root *mux.Router) error {
	// 注册认证模块
	auth := authn.New(user_repo.User(),
		authn.WithMiddleware(user_svc.Login().Middleware()),
	)
	authn.SetDefault(auth)

	r := root.Group("/api/v1")

	userLoginCtr := user_ctr.NewLogin()
	{
		// 绑定路由
		r.Group("/").Bind(
			userLoginCtr.Register,
			userLoginCtr.Login,
		)

		r.Group("/", auth.Middleware(true)).Bind(
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
	}

	return nil
}
