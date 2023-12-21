package mux

import (
	"context"
	"encoding/json"
	"github.com/codfrm/cago/pkg/utils"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"reflect"
	"strconv"
)

type bind struct {
	ctx *gin.Context
}

func ShouldBindWith(c *gin.Context, obj any) error {
	return c.ShouldBindWith(obj, &bind{ctx: c})
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
	var form func(key string) []string
	if req.Method == http.MethodGet ||
		req.Method == http.MethodDelete {
		form = b.ctx.QueryArray
	} else {
		switch b.ctx.ContentType() {
		case binding.MIMEJSON:
			if req == nil || req.Body == nil {
				return httputils.NewInternalServerError(-1, "json body is nil")
			}
			decoder := json.NewDecoder(req.Body)
			if err := decoder.Decode(ptr); err != nil {
				return httputils.NewInternalServerError(-1, err.Error())
			}
		case binding.MIMEMultipartPOSTForm, binding.MIMEPOSTForm:
			form = b.ctx.PostFormArray
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
				if err := b.bind(req, ptrElem.Field(i).Addr().Interface()); err != nil {
					return err
				}
			} else if form != nil {
				// 处理key,label的情况,例如: key,default=1
				key, opts := utils.Head(key, ",")
				opts, val := utils.Head(opts, "=")
				if opts == "default" && len(form(key)) == 0 {
					setValue(ptrElem.Field(i), tag, []string{val})
				} else {
					setValue(ptrElem.Field(i), tag, form(key))
				}
			}
		} else if uri := tag.Get("uri"); uri != "" {
			setValue(ptrElem.Field(i), tag, []string{b.ctx.Param(uri)})
		} else if header := tag.Get("header"); header != "" {
			setValue(ptrElem.Field(i), tag, []string{b.ctx.GetHeader(header)})
		}
	}
	return nil
}

// 设置值,暂时只支持基础类型
func setValue(field reflect.Value, tag reflect.StructTag, value []string) {
	if len(value) == 0 {
		return
	}
	switch field.Kind() {
	case reflect.String:
		field.SetString(value[0])
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		i, _ := strconv.ParseInt(value[0], 10, 64)
		field.SetInt(i)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		i, _ := strconv.ParseUint(value[0], 10, 64)
		field.SetUint(i)
	case reflect.Float32, reflect.Float64:
		i, _ := strconv.ParseFloat(value[0], 64)
		field.SetFloat(i)
	case reflect.Bool:
		i, _ := strconv.ParseBool(value[0])
		field.SetBool(i)
	case reflect.Slice:
		if field.Type() == reflect.TypeOf([]string{}) {
			field.Set(reflect.ValueOf(value))
		}
	case reflect.Map:
		// JSON解析
		_ = json.Unmarshal([]byte(value[0]), field.Addr().Interface())
	default:
		if field.Type() == reflect.TypeOf(primitive.ObjectID{}) {
			if id, err := primitive.ObjectIDFromHex(value[0]); err == nil {
				field.Set(reflect.ValueOf(id))
			}
		}
	}
}
