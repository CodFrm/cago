package cache

import (
	"context"
	"errors"
	"time"
)

var ErrDependNotValid = errors.New("depend not valid")

// KeyDepend key依赖
type KeyDepend struct {
	store Cache
	Key   string `json:"key"`
	Value int64  `json:"Value"`
}

func NewKeyDepend(store Cache, key string) *KeyDepend {
	return &KeyDepend{
		store: store,
		Key:   key,
	}
}

func WithKeyDepend(store Cache, key string) Option {
	return func(options *Options) {
		options.Depend = NewKeyDepend(store, key)
	}
}

// InvalidKey 使key失效
func (v *KeyDepend) InvalidKey(ctx context.Context) error {
	return v.store.Set(ctx, v.Key, &KeyDepend{Key: v.Key, Value: time.Now().Unix()}).Err()
}

// Val 获取依赖的值
func (v *KeyDepend) Val(ctx context.Context) interface{} {
	ret := &KeyDepend{}
	if err := v.store.Get(ctx, v.Key).Scan(ret); err != nil {
		if err := v.InvalidKey(ctx); err != nil {
			return err
		}
		return &KeyDepend{Key: v.Key, Value: time.Now().Unix()}
	}
	return ret
}

// Valid 检查依赖是否有效
func (v *KeyDepend) Valid(ctx context.Context) error {
	val := v.Val(ctx).(*KeyDepend)
	if v.Value == val.Value {
		return nil
	}
	return ErrDependNotValid
}

// NilDep 用于set的时候反序列化,减少一次dep判断
type NilDep struct {
	Depend
}

func (n *NilDep) Valid(ctx context.Context) error {
	return nil
}
