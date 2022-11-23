package cache

import (
	"context"
	"errors"
	"time"
)

type KeyDepend struct {
	store ICache
	Key   string `json:"key"`
	Value int64  `json:"value"`
}

func NewKeyDepend(store ICache, key string) *KeyDepend {
	return &KeyDepend{
		store: store,
		Key:   key,
	}
}

func WithKeyDepend(store ICache, key string) Option {
	return func(options *Options) {
		options.depend = NewKeyDepend(store, key)
	}
}

func (v *KeyDepend) InvalidKey(ctx context.Context) error {
	return v.store.Set(ctx, v.Key, &KeyDepend{Key: v.Key, Value: time.Now().Unix()})
}

func (v *KeyDepend) Val(ctx context.Context) interface{} {
	ret := &KeyDepend{}
	if err := v.store.Get(ctx, v.Key, ret); err != nil {
		if err := v.InvalidKey(ctx); err != nil {
			return err
		}
		return &KeyDepend{Key: v.Key, Value: time.Now().Unix()}
	}
	return ret
}

func (v *KeyDepend) Ok(ctx context.Context) error {
	val := v.Val(ctx).(*KeyDepend)
	if v.Value == val.Value {
		return nil
	}
	return errors.New("val not equal")
}
