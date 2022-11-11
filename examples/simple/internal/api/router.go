package api

import (
	"github.com/codfrm/cago/examples/simple/internal/controller/user"
	"github.com/codfrm/cago/server/http"
)

func Router(r *http.Router) error {
	return r.Group("/").Bind(
		user.NewUser(),
	)
}
