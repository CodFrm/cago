package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// CtxRedis 简化操作,慢慢封装
type CtxRedis interface {
	Get(key string) *redis.StringCmd
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(key string) *redis.IntCmd

	HGet(key, field string) *redis.StringCmd
	HSet(key string, value ...interface{}) *redis.IntCmd
	HDel(key string, fields ...string) *redis.IntCmd
}

type ctxRedis struct {
	client *redis.Client
	ctx    context.Context
}

func (c *ctxRedis) Get(key string) *redis.StringCmd {
	return c.client.Get(c.ctx, key)
}

func (c *ctxRedis) Del(key string) *redis.IntCmd {
	return c.client.Del(c.ctx, key)
}

func (c *ctxRedis) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return c.client.Set(c.ctx, key, value, expiration)
}

func (c *ctxRedis) HGet(key, field string) *redis.StringCmd {
	return c.client.HGet(c.ctx, key, field)
}

func (c *ctxRedis) HSet(key string, value ...interface{}) *redis.IntCmd {
	return c.client.HSet(c.ctx, key, value...)
}

func (c *ctxRedis) HDel(key string, fields ...string) *redis.IntCmd {
	return c.client.HDel(c.ctx, key, fields...)
}
