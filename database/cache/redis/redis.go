package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/codfrm/cago/database/cache/cache"
	"github.com/codfrm/cago/pkg/trace"
	"github.com/go-redis/redis/extra/redisotel/v9"
	"github.com/go-redis/redis/v9"
)

type redisCache struct {
	redis *redis.Client
}

func NewRedisCache(config *redis.Options) (cache.Cache, error) {
	client := redis.NewClient(config)
	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}
	if tp := trace.Default(); tp != nil {
		if err := redisotel.InstrumentTracing(client,
			redisotel.WithTracerProvider(tp),
			redisotel.WithDBSystem("cache"),
		); err != nil {
			return nil, err
		}
	}
	return &redisCache{
		redis: client,
	}, nil
}

type getOrSetValue struct {
	cache.Value
	set func() cache.Value
}

func (g *getOrSetValue) Scan(v interface{}) error {
	err := g.Value.Scan(v)
	if err != nil {
		if err == cache.ErrDependNotValid {
			return g.set().Scan(v)
		}
		return err
	}
	return nil
}

func (r *redisCache) GetOrSet(ctx context.Context, key string, set func() (interface{}, error), opts ...cache.Option) cache.Value {
	ret := r.Get(ctx, key, opts...)
	if ret.Err() != nil {
		val, err := set()
		if err != nil {
			return newValue(ctx, "", cache.NewOptions(opts...), err)
		}
		return r.Set(ctx, key, val, opts...)
	}
	return &getOrSetValue{Value: ret, set: func() cache.Value {
		val, err := set()
		if err != nil {
			return newValue(ctx, "", cache.NewOptions(opts...), err)
		}
		return r.Set(ctx, key, val, opts...)
	}}
}

func (r *redisCache) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (r *redisCache) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (r *redisCache) Get(ctx context.Context, key string, opts ...cache.Option) cache.Value {
	data, err := r.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		err = cache.ErrNotFound
	}
	options := cache.NewOptions(opts...)
	return newValue(ctx, data, options, err)
}

// 用于set的时候反序列化,减少一次dep判断
type nilDep struct {
	cache.Depend
}

func (n *nilDep) Valid(ctx context.Context) error {
	return nil
}

func (r *redisCache) Set(ctx context.Context, key string, val interface{}, opts ...cache.Option) cache.Value {
	options := cache.NewOptions(opts...)
	ttl := time.Duration(0)
	if options.Expiration > 0 {
		ttl = options.Expiration
	}
	data, err := Marshal(ctx, val, options)
	if err != nil {
		return newValue(ctx, "", options, err)
	}
	s := string(data)
	if err := r.redis.Set(ctx, key, s, ttl).Err(); err != nil {
		return newValue(ctx, "", options, err)
	}
	if options.Depend != nil {
		// 移除掉依赖
		options.Depend = &nilDep{}
	}
	return newValue(ctx, s, options, err)
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
