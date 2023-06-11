package api

import (
	cache2 "github.com/codfrm/cago/database/cache"
	_ "github.com/codfrm/cago/examples/simple/docs"
	"github.com/codfrm/cago/examples/simple/internal/controller/example_ctr"
	"github.com/codfrm/cago/examples/simple/internal/controller/user_ctr"
	"github.com/codfrm/cago/examples/simple/internal/repository/user_repo"
	"github.com/codfrm/cago/middleware/sessions"
	"github.com/codfrm/cago/middleware/sessions/cache"
	"github.com/codfrm/cago/server/mux"
)

// Router 路由
// @title    api文档
// @version  1.0
// @BasePath /api/v1
func Router(root *mux.Router) error {
	// 注册储存实例
	user_repo.RegisterUser(user_repo.NewUser())

	r := root.Group("/api/v1", sessions.Middleware("simple-session", cache.NewCacheStore(
		cache2.Default(), "simple",
	)))

	userLoginCtr := user_ctr.NewLogin()
	{
		// 绑定路由
		r.Group("/").Bind(
			userLoginCtr.Register,
			userLoginCtr.Login,
		)
		r.Group("/", userLoginCtr.Middleware()).Bind(
			userLoginCtr.Logout,
		)
	}

	{
		exampleCtl := example_ctr.NewExample()
		r.Group("/").Bind(
			exampleCtl.Ping,
		)
		r.Group("/", userLoginCtr.Middleware()).Bind(
			exampleCtl.Login,
		)
	}

	return nil
}
