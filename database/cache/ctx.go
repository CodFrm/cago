package cache

import (
	"context"

	"github.com/codfrm/cago/database/cache/cache"
)

type CtxCache struct {
	cache.Cache
	ctx context.Context
}

func Ctx(ctx context.Context) *CtxCache {
	return &CtxCache{ctx: ctx, Cache: Default()}
}

func (c *CtxCache) GetOrSet(key string, set func() (interface{}, error), opts ...cache.Option) cache.Value {
	return c.Cache.GetOrSet(c.ctx, key, set, opts...)
}

func (c *CtxCache) Set(key string, val interface{}, opts ...cache.Option) cache.Value {
	return c.Cache.Set(c.ctx, key, val, opts...)
}

func (c *CtxCache) Get(key string, opts ...cache.Option) cache.Value {
	return c.Cache.Get(c.ctx, key, opts...)
}

func (c *CtxCache) Has(key string) (bool, error) {
	return c.Cache.Has(c.ctx, key)
}

func (c *CtxCache) Del(key string) error {
	return c.Cache.Del(c.ctx, key)
}
