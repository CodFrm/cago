package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type value struct {
	ctx     context.Context
	data    string
	err     error
	options *Options
}

func NewValue(ctx context.Context, data string, options *Options, err ...error) Value {
	var e error
	if len(err) != 0 {
		e = err[0]
	}
	return &value{
		ctx:     ctx,
		data:    data,
		err:     e,
		options: options,
	}
}

func (v *value) Bytes() ([]byte, error) {
	return []byte(v.data), v.err
}

func (v *value) Result() (string, error) {
	return v.data, v.err
}

func (v *value) Int64() (int64, error) {
	if v.err != nil {
		return 0, v.err
	}
	return strconv.ParseInt(v.data, 10, 64)
}

func (v *value) Bool() (bool, error) {
	if v.err != nil {
		return false, v.err
	}
	return strconv.ParseBool(v.data)
}

func (v *value) Err() error {
	return v.err
}

func (v *value) Scan(data interface{}) error {
	if v.err != nil {
		return v.err
	}
	return Unmarshal(v.ctx, []byte(v.data), data, v.options)
}

// 带依赖的缓存数据
type dependStore struct {
	Depend interface{} `json:"depend"`
	Data   interface{} `json:"data"`
}

func Unmarshal(ctx context.Context, data []byte, v interface{}, options *Options) error {
	// 反序列化时,如果有依赖,带上依赖
	if options.Depend != nil {
		newV := reflect.New(reflect.TypeOf(v).Elem())
		dependStore := &dependStore{
			Depend: options.Depend,
			Data:   newV.Interface(),
		}
		if err := json.Unmarshal(data, dependStore); err != nil {
			return err
		}
		if err := options.Depend.Valid(ctx); err != nil {
			return err
		}
		// 设置值
		reflect.ValueOf(v).Elem().Set(newV.Elem())
		return nil
	} else {
		// 否则直接反序列化
		return json.Unmarshal(data, v)
	}
}

func Marshal(ctx context.Context, data interface{}, options *Options) ([]byte, error) {
	if options.Depend != nil {
		dependStore := &dependStore{
			Depend: options.Depend.Val(ctx),
			Data:   data,
		}
		return json.Marshal(dependStore)
	} else {
		// 基础类型直接转成字符串
		switch v := data.(type) {
		case string:
			return []byte(v), nil
		case []byte:
			return v, nil
		case int8, int16, int32, int64, int,
			uint8, uint16, uint32, uint64, uint,
			float32, float64, complex64, complex128:
			return []byte(fmt.Sprintf("%v", v)), nil
		case bool:
			return []byte(strconv.FormatBool(v)), nil
		default:
			return json.Marshal(data)
		}
	}
}

type GetOrSetValue struct {
	Value
	Set func() Value
}

func (g *GetOrSetValue) Scan(v interface{}) error {
	err := g.Value.Scan(v)
	if err != nil {
		if errors.Is(err, ErrDependNotValid) {
			return g.Set().Scan(v)
		}
		return err
	}
	return nil
}
