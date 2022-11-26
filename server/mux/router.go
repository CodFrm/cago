package mux

import (
	"context"
	"errors"
	"net/http"
	"reflect"

	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/gin-gonic/gin"
)

type Router struct {
	gin.IRouter
}

func (r *Router) Group(path string) *Router {
	return &Router{r.IRouter.Group(path)}
}

type RegisterRouter interface {
	Router(r *Router) error
}

// Bind 绑定控制器
func (r *Router) Bind(controller ...interface{}) error {
	// 反射解析控制器方法
	for _, c := range controller {
		el := reflect.TypeOf(c)
		if el.Kind() == reflect.Ptr {
			reg, ok := c.(RegisterRouter)
			if ok {
				if err := reg.Router(r); err != nil {
					return err
				}
			}
			for i := 0; i < el.NumMethod(); i++ {
				method := el.Method(i)
				if err := r.bindFunc(reflect.ValueOf(c), method.Func); err != nil {
					return err
				}
			}
		} else if el.Kind() == reflect.Func {
			if err := r.bindFunc(reflect.Zero(el), reflect.ValueOf(c)); err != nil {
				return err
			}
		}
	}
	return nil
}

// 根据方法去绑定路由
func (r *Router) bindFunc(controller reflect.Value, method reflect.Value) error {
	methodType := method.Type()
	// 判断返回值是否是error
	if methodType.NumOut() != 2 || methodType.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
		return nil
	}
	// 获取路由参数
	request := methodType.In(2).Elem()
	route, ok := request.FieldByName("Meta")
	// 必须有Route字段
	if !ok || route.Type != reflect.TypeOf(Meta{}) {
		return errors.New("invalid method, second parameter must have Meta field")
	}
	ginContext := false
	if route.Tag.Get("context") == "gin" {
		ginContext = true
	}
	// 判断方法的第一个参数是否是context.Context
	parame1 := methodType.In(1)
	if methodType.NumIn() != 3 ||
		(parame1 != reflect.TypeOf((*context.Context)(nil)).Elem() &&
			parame1 != reflect.TypeOf((*gin.Context)(nil))) {
		return errors.New("invalid method, first parameter must be context.Context or *gin.Context")
	}
	var call func(a reflect.Value, b interface{}) []reflect.Value
	if controller.IsValid() {
		call = func(a reflect.Value, b interface{}) []reflect.Value {
			return method.Call([]reflect.Value{controller, a, reflect.ValueOf(b)})
		}
	} else {
		call = func(a reflect.Value, b interface{}) []reflect.Value {
			return method.Call([]reflect.Value{a, reflect.ValueOf(b)})
		}
	}
	switch route.Tag.Get("method") {
	case http.MethodGet:
		r.GET(route.Tag.Get("path"), r.bindHandler(request, call, ginContext))
	case http.MethodPost:
		r.POST(route.Tag.Get("path"), r.bindHandler(request, call, ginContext))
	case http.MethodPut:
		r.PUT(route.Tag.Get("path"), r.bindHandler(request, call, ginContext))
	case http.MethodDelete:
		r.DELETE(route.Tag.Get("path"), r.bindHandler(request, call, ginContext))
	}
	return nil
}

func (r *Router) bindHandler(request reflect.Type, call func(a reflect.Value, b interface{}) []reflect.Value, ginContext bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建请求参数
		req := reflect.New(request)
		// 绑定请求参数
		i := req.Interface()
		if err := c.ShouldBind(i); err != nil {
			httputils.HandleResp(c, err)
			return
		}
		// 获取uri参数
		if err := c.ShouldBindUri(i); err != nil {
			httputils.HandleResp(c, err)
			return
		}
		// 调用控制器方法
		var ctx reflect.Value
		if ginContext {
			ctx = reflect.ValueOf(c)
		} else {
			ctx = reflect.ValueOf(c.Request.Context())
		}
		resp := call(ctx, req.Interface())
		if resp[1].IsNil() {
			httputils.HandleResp(c, resp[0].Interface())
			return
		}
		httputils.HandleResp(c, resp[1].Interface())
	}
}
