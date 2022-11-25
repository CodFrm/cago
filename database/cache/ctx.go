package cache

import (
	"context"

	"github.com/codfrm/cago/database/cache/cache"
)

type CtxCache struct {
	ctx   context.Context
	cache cache.Cache
}

func Ctx() *CtxCache {
	return &CtxCache{}
}

func (c *CtxCache) GetOrSet(key string, set func() (interface{}, error), opts ...cache.Option) cache.Value {
	return c.cache.GetOrSet(c.ctx, key, set, opts...)
}

func (c *CtxCache) Set(ctx context.Context, key string, val interface{}, opts ...cache.Option) cache.Value {
	return c.cache.Set(c.ctx, key, val, opts...)
}

func (c *CtxCache) Get(ctx context.Context, key string, get interface{}, opts ...cache.Option) cache.Value {
	return c.cache.Get(c.ctx, key, opts...)
}

func (c *CtxCache) Has(ctx context.Context, key string) (bool, error) {
	return c.cache.Has(c.ctx, key)
}

func (c *CtxCache) Del(ctx context.Context, key string) error {
	return c.cache.Del(c.ctx, key)
}
