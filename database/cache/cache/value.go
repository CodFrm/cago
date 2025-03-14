package cache

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
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
	var val []byte
	if err := v.Scan(&val); err != nil {
		return nil, err
	}
	return val, nil
}

func (v *value) Result() (string, error) {
	var val string
	if err := v.Scan(&val); err != nil {
		return "", err
	}
	return val, nil
}

func (v *value) Int64() (int64, error) {
	if v.err != nil {
		return 0, v.err
	}
	var val int64
	if err := v.Scan(&val); err != nil {
		return 0, err
	}
	return val, nil
}

func (v *value) Bool() (bool, error) {
	if v.err != nil {
		return false, v.err
	}
	var val bool
	if err := v.Scan(&val); err != nil {
		return false, err
	}
	return val, nil
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
		dependValue, err := options.Depend.ValInterface()
		if err != nil {
			return err
		}
		dependStore := &dependStore{
			Depend: dependValue,
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
		val, err := options.Depend.Val(ctx)
		if err != nil {
			return nil, err
		}
		dependStore := &dependStore{
			Depend: val,
			Data:   data,
		}
		return json.Marshal(dependStore)
	} else {
		// 直接序列化
		return json.Marshal(data)
	}
}

type GetOrSetValue struct {
	Value
	Set func() Value
}

func (g *GetOrSetValue) Int64() (int64, error) {
	var val int64
	if err := g.Scan(&val); err != nil {
		return 0, err
	}
	return val, nil
}

func (g *GetOrSetValue) Bool() (bool, error) {
	var val bool
	if err := g.Scan(&val); err != nil {
		return false, err
	}
	return val, nil
}

func (g *GetOrSetValue) Bytes() ([]byte, error) {
	var val []byte
	if err := g.Scan(&val); err != nil {
		return nil, err
	}
	return val, nil
}

func (g *GetOrSetValue) Result() (string, error) {
	var val string
	if err := g.Scan(&val); err != nil {
		return "", err
	}
	return val, nil
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
