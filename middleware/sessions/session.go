package sessions

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const DefaultKey = "cago/session"

// Middleware 在gin-contrib/sessions上做封装
func Middleware(name string, store sessions.Store) gin.HandlerFunc {
	return sessions.Sessions(name, store)
}

func Ctx(ctx *gin.Context) sessions.Session {
	return sessions.Default(ctx)
}
