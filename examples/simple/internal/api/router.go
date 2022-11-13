package api

import (
	"github.com/codfrm/cago/examples/simple/internal/controller/user"
	"github.com/codfrm/cago/server/http"
)

// Router 路由
// @title    api文档
// @version  1.0
// @BasePath /api/v1
func Router(r *http.Router) error {
	return r.Group("/").Bind(
		user.NewUser(),
	)
}
