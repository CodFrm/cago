package cache

import (
	"context"
	"errors"
)

var (
	ErrNil = errors.New("cache: nil")
)

type Cache interface {
	GetOrSet(ctx context.Context, key string, set func() (interface{}, error), opts ...Option) Value
	Set(ctx context.Context, key string, val interface{}, opts ...Option) Value
	Get(ctx context.Context, key string, opts ...Option) Value
	Has(ctx context.Context, key string) (bool, error)
	Del(ctx context.Context, key string) error
	Close() error
}

type DependValue interface {
	Equate(DependValue) bool
}

type Depend interface {
	Val(ctx context.Context) (DependValue, error)
	ValInterface() (DependValue, error)
	Valid(ctx context.Context) error
}

type Value interface {
	Result() (string, error)
	Err() error
	Scan(v interface{}) error
	Bytes() ([]byte, error)
	Int64() (int64, error)
	Bool() (bool, error)
}

func IsNil(err error) bool {
	return errors.Is(err, ErrNil)
}

type prefixCache struct {
	Cache  Cache
	prefix string
}

func NewPrefixCache(prefix string, cache Cache) Cache {
	return &prefixCache{prefix: prefix, Cache: cache}
}

func (p *prefixCache) key(key string) string {
	return p.prefix + key
}

func (p *prefixCache) GetOrSet(ctx context.Context, key string, set func() (interface{}, error), opts ...Option) Value {
	return p.Cache.GetOrSet(ctx, p.key(key), set, opts...)
}

func (p *prefixCache) Set(ctx context.Context, key string, val interface{}, opts ...Option) Value {
	return p.Cache.Set(ctx, p.key(key), val, opts...)
}

func (p *prefixCache) Get(ctx context.Context, key string, opts ...Option) Value {
	return p.Cache.Get(ctx, p.key(key), opts...)
}

func (p *prefixCache) Has(ctx context.Context, key string) (bool, error) {
	return p.Cache.Has(ctx, p.prefix+key)
}

func (p *prefixCache) Del(ctx context.Context, key string) error {
	return p.Cache.Del(ctx, p.prefix+key)
}

func (p *prefixCache) Close() error {
	return p.Cache.Close()
}
