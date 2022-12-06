package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"
)

// CtxRedis 简化操作,慢慢封装
type CtxRedis struct {
	client *redis.Client
	ctx    context.Context
}

func (c *CtxRedis) Get(key string) *redis.StringCmd {
	return c.client.Get(c.ctx, key)
}

func (c *CtxRedis) Del(key string) *redis.IntCmd {
	return c.client.Del(c.ctx, key)
}

func (c *CtxRedis) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return c.client.Set(c.ctx, key, value, expiration)
}

func (c *CtxRedis) HGet(key, field string) *redis.StringCmd {
	return c.client.HGet(c.ctx, key, field)
}

func (c *CtxRedis) HSet(key string, value ...interface{}) *redis.IntCmd {
	return c.client.HSet(c.ctx, key, value...)
}

func (c *CtxRedis) HDel(key string, fields ...string) *redis.IntCmd {
	return c.client.HDel(c.ctx, key, fields...)
}
