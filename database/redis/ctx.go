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

func (c *CtxRedis) SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return c.client.SetNX(c.ctx, key, value, expiration)
}

func (c *CtxRedis) SetXX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return c.client.SetXX(c.ctx, key, value, expiration)
}

func (c *CtxRedis) SetEx(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return c.client.SetEx(c.ctx, key, value, expiration)
}

func (c *CtxRedis) SetRange(key string, offset int64, value string) *redis.IntCmd {
	return c.client.SetRange(c.ctx, key, offset, value)
}

func (c *CtxRedis) GetRange(key string, start, end int64) *redis.StringCmd {
	return c.client.GetRange(c.ctx, key, start, end)
}

func (c *CtxRedis) GetSet(key string, value interface{}) *redis.StringCmd {
	return c.client.GetSet(c.ctx, key, value)
}

func (c *CtxRedis) Incr(key string) *redis.IntCmd {
	return c.client.Incr(c.ctx, key)
}

func (c *CtxRedis) IncrBy(key string, value int64) *redis.IntCmd {
	return c.client.IncrBy(c.ctx, key, value)
}

func (c *CtxRedis) IncrByFloat(key string, value float64) *redis.FloatCmd {
	return c.client.IncrByFloat(c.ctx, key, value)
}

func (c *CtxRedis) Decr(key string) *redis.IntCmd {
	return c.client.Decr(c.ctx, key)
}

func (c *CtxRedis) DecrBy(key string, decrement int64) *redis.IntCmd {
	return c.client.DecrBy(c.ctx, key, decrement)
}

func (c *CtxRedis) Exists(keys ...string) *redis.IntCmd {
	return c.client.Exists(c.ctx, keys...)
}

func (c *CtxRedis) Expire(key string, expiration time.Duration) *redis.BoolCmd {
	return c.client.Expire(c.ctx, key, expiration)
}

func (c *CtxRedis) ExpireAt(key string, tm time.Time) *redis.BoolCmd {
	return c.client.ExpireAt(c.ctx, key, tm)
}

func (c *CtxRedis) TTL(key string) *redis.DurationCmd {
	return c.client.TTL(c.ctx, key)
}

func (c *CtxRedis) LPush(key string, values ...interface{}) *redis.IntCmd {
	return c.client.LPush(c.ctx, key, values...)
}

func (c *CtxRedis) RPush(key string, values ...interface{}) *redis.IntCmd {
	return c.client.RPush(c.ctx, key, values...)
}

func (c *CtxRedis) LPop(key string) *redis.StringCmd {
	return c.client.LPop(c.ctx, key)
}

func (c *CtxRedis) RPop(key string) *redis.StringCmd {
	return c.client.RPop(c.ctx, key)
}

func (c *CtxRedis) LLen(key string) *redis.IntCmd {
	return c.client.LLen(c.ctx, key)
}

func (c *CtxRedis) LRange(key string, start, stop int64) *redis.StringSliceCmd {
	return c.client.LRange(c.ctx, key, start, stop)
}

func (c *CtxRedis) LTrim(key string, start, stop int64) *redis.StatusCmd {
	return c.client.LTrim(c.ctx, key, start, stop)
}

func (c *CtxRedis) LRem(key string, count int64, value interface{}) *redis.IntCmd {
	return c.client.LRem(c.ctx, key, count, value)
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

func (c *CtxRedis) HGetAll(key string) *redis.MapStringStringCmd {
	return c.client.HGetAll(c.ctx, key)
}

func (c *CtxRedis) HIncrBy(key, field string, incr int64) *redis.IntCmd {
	return c.client.HIncrBy(c.ctx, key, field, incr)
}

func (c *CtxRedis) HIncrByFloat(key, field string, incr float64) *redis.FloatCmd {
	return c.client.HIncrByFloat(c.ctx, key, field, incr)
}

func (c *CtxRedis) HExists(key, field string) *redis.BoolCmd {
	return c.client.HExists(c.ctx, key, field)
}

func (c *CtxRedis) HKeys(key string) *redis.StringSliceCmd {
	return c.client.HKeys(c.ctx, key)
}

func (c *CtxRedis) HLen(key string) *redis.IntCmd {
	return c.client.HLen(c.ctx, key)
}

func (c *CtxRedis) HSetNX(key, field string, value interface{}) *redis.BoolCmd {
	return c.client.HSetNX(c.ctx, key, field, value)
}

func (c *CtxRedis) HVals(key string) *redis.StringSliceCmd {
	return c.client.HVals(c.ctx, key)
}

func (c *CtxRedis) HScan(key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return c.client.HScan(c.ctx, key, cursor, match, count)
}

func (c *CtxRedis) PFCount(keys ...string) *redis.IntCmd {
	return c.client.PFCount(c.ctx, keys...)
}

func (c *CtxRedis) PFAdd(key string, els ...interface{}) *redis.IntCmd {
	return c.client.PFAdd(c.ctx, key, els...)
}

func (c *CtxRedis) PFMerge(dest string, keys ...string) *redis.StatusCmd {
	return c.client.PFMerge(c.ctx, dest, keys...)
}

func (c *CtxRedis) ScanType(cursor uint64, match string, count int64, keyType string) *redis.ScanCmd {
	return c.client.ScanType(c.ctx, cursor, match, count, keyType)
}
