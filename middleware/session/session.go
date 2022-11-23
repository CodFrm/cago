package session

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

const DefaultKey = "cago/session"

func Middleware(name string, store sessions.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := store.New(c.Request, name)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Set(DefaultKey, session)
	}
}

func Ctx(ctx *gin.Context) sessions.Store {
	v, _ := ctx.Get(DefaultKey)
	return v.(sessions.Store)
}
