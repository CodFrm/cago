package mux

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type bind struct {
	ctx *gin.Context
}

func (b *bind) Name() string {
	return "cago"
}

// Validate 数据校验
type Validate interface {
	Validate(ctx context.Context) error
}

func (b *bind) Bind(req *http.Request, ptr any) error {
	if err := b.bind(req, ptr); err != nil {
		return err
	}
	return binding.Validator.ValidateStruct(ptr)
}

func (b *bind) bind(req *http.Request, ptr any) error {
	// 根据tag绑定数据
	if v, ok := ptr.(Validate); ok {
		if err := v.Validate(b.ctx); err != nil {
			return err
		}
	}
	// Check if ptr is a map
	ptrVal := reflect.ValueOf(ptr)
	ptrElem := ptrVal.Elem()
	ptrType := ptrElem.Type()
	var form func(key string) string
	if req.Method == http.MethodGet ||
		req.Method == http.MethodDelete {
		form = b.ctx.Query
	} else {
		switch b.ctx.ContentType() {
		case binding.MIMEJSON:
			if req == nil || req.Body == nil {
				return errors.New("invalid request")
			}
			decoder := json.NewDecoder(req.Body)
			if err := decoder.Decode(ptr); err != nil {
				return err
			}
		case binding.MIMEMultipartPOSTForm, binding.MIMEPOSTForm:
			form = b.ctx.PostForm
		}
	}
	for i := 0; i < ptrElem.NumField(); i++ {
		tag := ptrType.Field(i).Tag
		if tag == "" {
			continue
		}
		if key := tag.Get("form"); key != "" {
			if key == "-" {
				continue
			}
			if key == ",inline" {
				b.bind(req, ptrElem.Field(i).Addr().Interface())
			} else if form != nil {
				// 处理key,label的情况,例如: key,default=1
				key, opts := head(key, ",")
				opts, val := head(opts, "=")
				if opts == "default" && form(key) == "" {
					setValue(ptrElem.Field(i), tag, val)
				} else {
					setValue(ptrElem.Field(i), tag, form(key))
				}
			}
		} else if uri := tag.Get("uri"); uri != "" {
			setValue(ptrElem.Field(i), tag, b.ctx.Param(uri))
		} else if header := tag.Get("header"); header != "" {
			setValue(ptrElem.Field(i), tag, b.ctx.GetHeader(header))
		}
	}
	return nil
}

// 设置值,暂时只支持基础类型
func setValue(field reflect.Value, tag reflect.StructTag, value string) {
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

func head(str, sep string) (head string, tail string) {
	idx := strings.Index(str, sep)
	if idx < 0 {
		return str, ""
	}
	return str[:idx], str[idx+len(sep):]
}
