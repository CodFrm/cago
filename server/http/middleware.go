package http

import (
	"context"

	"github.com/gin-gonic/gin"
)

type ginKVContext struct {
	context.Context
	ginCtx *gin.Context
}

func (g *ginKVContext) Value(key interface{}) interface{} {
	s, ok := key.(string)
	if !ok {
		return g.Context.Value(key)
	}
	value, exist := g.ginCtx.Get(s)
	if exist {
		return value
	}
	return g.Context.Value(key)
}

func GinKVContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = c.Request.WithContext(&ginKVContext{
			ginCtx:  c,
			Context: c.Request.Context(),
		})
	}
}
