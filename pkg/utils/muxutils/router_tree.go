package muxutils

import (
	"github.com/codfrm/cago/server/mux"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Router struct {
	Method       string
	RelativePath string
	Handler      gin.HandlerFunc
}

type RouterTree struct {
	Middleware []gin.HandlerFunc
	Handler    []interface{}
}

func BindTree(r *mux.Router, tree []*RouterTree) {
	for _, v := range tree {
		if len(v.Handler) > 0 {
			rg := r.Group("/")
			rg.Use(v.Middleware...)
			for _, handler := range v.Handler {
				switch h := handler.(type) {
				case *RouterTree:
					BindTree(rg, []*RouterTree{h})
				case *Router:
					switch h.Method {
					case http.MethodGet:
						rg.GET(h.RelativePath, h.Handler)
					case http.MethodPost:
						rg.POST(h.RelativePath, h.Handler)
					case http.MethodPut:
						rg.PUT(h.RelativePath, h.Handler)
					case http.MethodDelete:
						rg.DELETE(h.RelativePath, h.Handler)
					}
				default:
					rg.Bind(h)
				}
			}
		}
	}
}
