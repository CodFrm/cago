package cache

import (
	"context"
	"encoding/json"
	"reflect"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisCache struct {
	redis *redis.Client
}

func newRedisCache(redis *redis.Client) ICache {
	return &redisCache{
		redis: redis,
	}
}

func (r *redisCache) GetOrSet(ctx context.Context, key string, get interface{}, set func() (interface{}, error), opts ...Option) error {
	err := r.Get(ctx, key, get, opts...)
	if err != nil {
		val, err := set()
		if err != nil {
			return err
		}
		if err := r.Set(ctx, key, val, opts...); err != nil {
			return err
		}
		copyInterface(get, val)
	}
	return nil
}

func (r *redisCache) Get(ctx context.Context, key string, get interface{}, opts ...Option) error {
	val, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	options := NewOptions(opts...)
	ret := &data{Value: get, Depend: options.depend}
	if err := json.Unmarshal([]byte(val), ret); err != nil {
		return err
	}
	if options.depend != nil {
		if err := options.depend.Ok(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (r *redisCache) Set(ctx context.Context, key string, val interface{}, opts ...Option) error {
	options := NewOptions(opts...)
	ttl := time.Duration(0)
	if options.expiration > 0 {
		ttl = options.expiration
	}
	data := &data{Value: val}
	if options.depend != nil {
		data.Depend = options.depend.Val(ctx)
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if err := r.redis.Set(ctx, key, b, ttl).Err(); err != nil {
		return err
	}
	return nil
}

func (r *redisCache) Has(ctx context.Context, key string) (bool, error) {
	ok, err := r.redis.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if ok == 1 {
		return true, nil
	}
	return false, nil
}

func (r *redisCache) Del(ctx context.Context, key string) error {
	return r.redis.Del(ctx, key).Err()
}

func copyInterface(dst interface{}, src interface{}) {
	dstof := reflect.ValueOf(dst)
	if dstof.Kind() == reflect.Ptr {
		el := dstof.Elem()
		srcof := reflect.ValueOf(src)
		if srcof.Kind() == reflect.Ptr {
			el.Set(srcof.Elem())
		} else if src == nil {
			dst = nil
		} else {
			el.Set(srcof)
		}
	}
}
