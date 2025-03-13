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
	Key   string           `json:"key"`
	Value Int64DependValue `json:"Value"`
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
	return v.store.Set(ctx, v.Key, &KeyDepend{Key: v.Key, Value: Int64DependValue(time.Now().Unix())}).Err()
}

type Int64DependValue int64

func (i Int64DependValue) Equate(d DependValue) bool {
	return i == d.(Int64DependValue)
}

// Val 获取依赖的值
func (v *KeyDepend) Val(ctx context.Context) (DependValue, error) {
	var i int64
	if err := v.store.Get(ctx, v.Key).Scan(&i); err != nil {
		if err := v.InvalidKey(ctx); err != nil {
			return nil, err
		}
		return Int64DependValue(time.Now().Unix()), nil
	}
	return Int64DependValue(i), nil
}

// Valid 检查依赖是否有效
func (v *KeyDepend) Valid(ctx context.Context) error {
	val, err := v.Val(ctx)
	if err != nil {
		return err
	}
	if val.Equate(v.Value) {
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
