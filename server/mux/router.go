package mux

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/gin-gonic/gin"
)

type Router struct {
	gin.IRouter
}

func (r *Router) Group(path string, handler ...gin.HandlerFunc) *Router {
	return &Router{r.IRouter.Group(path, handler...)}
}

// ShouldBindWith 绑定控制器
func (r *Router) Bind(handler ...interface{}) {
	// 反射解析控制器方法
	for _, c := range handler {
		el := reflect.TypeOf(c)
		if el.Kind() == reflect.Func {
			if err := r.bindFunc(reflect.Zero(el), reflect.ValueOf(c), true); err != nil {
				panic(err)
			}
		} else {
			panic("invalid controller")
		}
	}
}

// 根据方法去绑定路由
func (r *Router) bindFunc(controller reflect.Value, method reflect.Value, isFunc bool) error {
	methodType := method.Type()
	// 判断返回值是否是error
	if methodType.NumOut() != 2 || methodType.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
		return nil
	}
	pos := 0
	if isFunc {
		pos = -1
	}
	// 获取路由参数
	request := methodType.In(2 + pos).Elem()
	route, ok := request.FieldByName("Meta")
	// 必须有Route字段
	if !ok || route.Type != reflect.TypeOf(Meta{}) {
		return errors.New("invalid method, second parameter must have Meta field")
	}
	ginContext := false
	// 判断方法的第一个参数是否是context.Context
	parame1 := methodType.In(1 + pos)
	if methodType.NumIn() != 3+pos {
		return errors.New("invalid method, first parameter must be context.Context or *gin.Context")
	}
	if parame1 != reflect.TypeOf((*context.Context)(nil)).Elem() {
		if parame1 != reflect.TypeOf((*gin.Context)(nil)) {
			return errors.New("invalid method, first parameter must be context.Context or *gin.Context")
		}
		ginContext = true
	}
	var call func(a reflect.Value, b interface{}) []reflect.Value
	if !isFunc {
		call = func(a reflect.Value, b interface{}) []reflect.Value {
			return method.Call([]reflect.Value{controller, a, reflect.ValueOf(b)})
		}
	} else {
		call = func(a reflect.Value, b interface{}) []reflect.Value {
			return method.Call([]reflect.Value{a, reflect.ValueOf(b)})
		}
	}

	methods := strings.Split(route.Tag.Get("method"), ",")
	for _, method := range methods {
		switch method {
		case http.MethodPost:
			r.POST(route.Tag.Get("path"), r.bindHandler(request, call, ginContext))
		case http.MethodPut:
			r.PUT(route.Tag.Get("path"), r.bindHandler(request, call, ginContext))
		case http.MethodDelete:
			r.DELETE(route.Tag.Get("path"), r.bindHandler(request, call, ginContext))
		case http.MethodOptions:
			r.OPTIONS(route.Tag.Get("path"), r.bindHandler(request, call, ginContext))
		default:
			r.GET(route.Tag.Get("path"), r.bindHandler(request, call, ginContext))
		}
	}
	return nil
}

func (r *Router) bindHandler(request reflect.Type, call func(a reflect.Value, b interface{}) []reflect.Value, ginContext bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建请求参数
		req := reflect.New(request)
		// 绑定请求参数
		i := req.Interface()
		if err := ShouldBindWith(c, i); err != nil {
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
