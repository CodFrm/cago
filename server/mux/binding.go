package mux

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type bind struct {
	ctx *gin.Context
}

func (b *bind) Name() string {
	return "cago"
}

func (b *bind) Bind(req *http.Request, ptr any) error {
	// 根据tag绑定数据
	// Check if ptr is a map
	ptrVal := reflect.ValueOf(ptr)
	ptrElem := ptrVal.Elem()
	ptrType := ptrElem.Type()
	flag := true
	switch b.ctx.ContentType() {
	case binding.MIMEJSON:
		if req == nil || req.Body == nil {
			return errors.New("invalid request")
		}
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(ptr); err != nil {
			return err
		}
		flag = false
	}
	for i := 0; i < ptrElem.NumField(); i++ {
		tag := ptrType.Field(i).Tag
		if form := tag.Get("form"); form != "" && flag {
			setValue(ptrElem.Field(i), b.ctx.PostForm(form))
		} else if uri := tag.Get("uri"); uri != "" {
			setValue(ptrElem.Field(i), b.ctx.Param(uri))
		}
	}
	return binding.Validator.ValidateStruct(ptr)
}

// 设置值,暂时只支持基础类型
func setValue(field reflect.Value, value string) {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		i, _ := strconv.ParseInt(value, 10, 64)
		field.SetInt(i)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		i, _ := strconv.ParseUint(value, 10, 64)
		field.SetUint(i)
	case reflect.Float32, reflect.Float64:
		i, _ := strconv.ParseFloat(value, 64)
		field.SetFloat(i)
	case reflect.Bool:
		i, _ := strconv.ParseBool(value)
		field.SetBool(i)
	}
}
