package http

import (
	"context"
	"errors"
	"net/http"
	"reflect"

	"github.com/codfrm/cago/pkg/httputils"
	"github.com/gin-gonic/gin"
)

type Router struct {
	gin.IRouter
}

func (r *Router) Group(path string) *Router {
	return &Router{r.IRouter.Group(path)}
}

// Bind 绑定控制器
func (r *Router) Bind(controller ...interface{}) error {
	// 反射解析控制器方法
	for _, c := range controller {
		el := reflect.TypeOf(c)
		for i := 0; i < el.NumMethod(); i++ {
			method := el.Method(i)
			if err := r.bindMethod(reflect.ValueOf(c), method); err != nil {
				return err
			}
		}
	}
	return nil
}

// 根据方法去绑定路由
func (r *Router) bindMethod(controller reflect.Value, method reflect.Method) error {
	methodType := method.Func.Type()
	// 判断方法的第一个参数是否是context.Context
	if methodType.NumIn() != 3 || methodType.In(1) != reflect.TypeOf((*context.Context)(nil)).Elem() {
		return errors.New("invalid method, first parameter must be context.Context")
	}
	// 判断返回值是否是error
	if methodType.NumOut() != 2 || methodType.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
		return errors.New("invalid method, second return value must be error")
	}
	// 获取路由参数
	request := methodType.In(2).Elem()
	route, ok := request.FieldByName("Route")
	// 必须有Route字段
	if !ok || route.Type != reflect.TypeOf(Route{}) {
		return errors.New("invalid method, second parameter must have Route field")
	}
	switch route.Tag.Get("method") {
	case http.MethodGet:
		r.GET(route.Tag.Get("path"), r.bindHandler(controller, method, request))
	case http.MethodPost:
		r.POST(route.Tag.Get("path"), r.bindHandler(controller, method, request))
	case http.MethodPut:
		r.PUT(route.Tag.Get("path"), r.bindHandler(controller, method, request))
	case http.MethodDelete:
		r.DELETE(route.Tag.Get("path"), r.bindHandler(controller, method, request))
	}
	return nil
}

func (r *Router) bindHandler(controller reflect.Value, method reflect.Method, request reflect.Type) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建请求参数
		req := reflect.New(request)
		// 绑定请求参数
		i := req.Interface()
		if err := c.ShouldBind(i); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// 调用控制器方法
		resp := method.Func.Call([]reflect.Value{
			controller,
			reflect.ValueOf(c.Request.Context()),
			reflect.ValueOf(i),
		})
		if resp[1].IsNil() {
			httputils.HandleResp(c, resp[0].Interface())
			return
		}
		httputils.HandleResp(c, resp[1].Interface())
	}
}
