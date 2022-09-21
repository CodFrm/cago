package mux

import (
	"github.com/codfrm/cago/pkg/utils/validator"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
)

type HandlerFunc func(c *Context)

type Mux struct {
	engine  *gin.Engine
	group   *RouterGroup
	options *Options
}

func New(logger *zap.Logger, opts ...Option) *Mux {
	options := &Options{}
	for _, o := range opts {
		o(options)
	}
	binding.Validator = validator.NewValidator()
	e := gin.New()
	if options.tracerProvider != nil {
		e.Use(otelgin.Middleware(options.serviceName, otelgin.WithTracerProvider(options.tracerProvider)))
	}
	group := &RouterGroup{
		group:  e.Group(""),
		logger: logger,
	}
	group.group.Use(initHandler)
	return &Mux{
		engine: e,
		group:  group,
	}
}

// Run 通过addr启动http服务
func (m *Mux) Run(addr ...string) error {
	return m.engine.Run(addr...)
}

// Group 返回路由组
func (m *Mux) Group() *RouterGroup {
	return m.group
}

type RouterGroup struct {
	group  *gin.RouterGroup
	logger *zap.Logger
}

func (r *RouterGroup) Use(handlers ...HandlerFunc) *RouterGroup {
	r.group.Use(r.WrapHandlers(handlers...)...)
	return r
}

func (r *RouterGroup) Group(relativePath string, handlers ...HandlerFunc) *RouterGroup {
	group := r.group.Group(relativePath, r.WrapHandlers(handlers...)...)
	return &RouterGroup{group: group}
}

func (r *RouterGroup) Any(relativePath string, handlers ...HandlerFunc) {
	r.group.Any(relativePath, r.WrapHandlers(handlers...)...)
}

func (r *RouterGroup) GET(relativePath string, handlers ...HandlerFunc) {
	r.group.GET(relativePath, r.WrapHandlers(handlers...)...)
}

func (r *RouterGroup) POST(relativePath string, handlers ...HandlerFunc) {
	r.group.POST(relativePath, r.WrapHandlers(handlers...)...)
}

func (r *RouterGroup) DELETE(relativePath string, handlers ...HandlerFunc) {
	r.group.DELETE(relativePath, r.WrapHandlers(handlers...)...)
}

func (r *RouterGroup) PATCH(relativePath string, handlers ...HandlerFunc) {
	r.group.PATCH(relativePath, r.WrapHandlers(handlers...)...)
}

func (r *RouterGroup) PUT(relativePath string, handlers ...HandlerFunc) {
	r.group.PUT(relativePath, r.WrapHandlers(handlers...)...)
}

func (r *RouterGroup) OPTIONS(relativePath string, handlers ...HandlerFunc) {
	r.group.OPTIONS(relativePath, r.WrapHandlers(handlers...)...)
}

func (r *RouterGroup) HEAD(relativePath string, handlers ...HandlerFunc) {
	r.group.HEAD(relativePath, r.WrapHandlers(handlers...)...)
}

func (r *RouterGroup) WrapHandlers(handlers ...HandlerFunc) []gin.HandlerFunc {
	funcs := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		handler := handler
		funcs[i] = func(c *gin.Context) {
			logger, _ := c.Get("logger")
			handler(NewContext(c, logger.(*zap.Logger)))
		}
	}
	return funcs
}
