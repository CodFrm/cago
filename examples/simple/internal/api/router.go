package api

import (
	"github.com/codfrm/cago/examples/simple/internal/controller/user_ctr"
	"github.com/codfrm/cago/examples/simple/internal/repository"
	"github.com/codfrm/cago/examples/simple/internal/repository/persistence"
	"github.com/codfrm/cago/server/mux"
)

// Router 路由
// @title    api文档
// @version  1.0
// @BasePath /api/v1
func Router(r *mux.Router) error {
	repository.RegisterUser(persistence.NewUser())

	user := user_ctr.NewUser()
	r.Group("/").Bind(
		user.Register,
		user.Login,
		user.Logout,
	)
	return nil
}
