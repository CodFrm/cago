package cache

import "context"

type CtxCache struct {
	ctx   context.Context
	cache ICache
}

func Ctx() *CtxCache {
	return &CtxCache{}
}

func (c *CtxCache) GetOrSet(key string, get interface{}, set func() (interface{}, error), opts ...Option) error {
	return c.cache.GetOrSet(c.ctx, key, get, set, opts...)
}

func (c *CtxCache) Set(ctx context.Context, key string, val interface{}, opts ...Option) error {
	return c.cache.Set(c.ctx, key, val, opts...)
}

func (c *CtxCache) Get(ctx context.Context, key string, get interface{}, opts ...Option) error {
	return c.cache.Get(c.ctx, key, get, opts...)
}

func (c *CtxCache) Has(ctx context.Context, key string) (bool, error) {
	return c.cache.Has(c.ctx, key)
}

func (c *CtxCache) Del(ctx context.Context, key string) error {
	return c.cache.Del(c.ctx, key)
}
