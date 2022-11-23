package cache

import (
	"context"

	"github.com/codfrm/cago/configs"
	"github.com/go-redis/redis/v8"
)

type ICache interface {
	GetOrSet(ctx context.Context, key string, get interface{}, set func() (interface{}, error), opts ...Option) error
	Set(ctx context.Context, key string, val interface{}, opts ...Option) error
	Get(ctx context.Context, key string, get interface{}, opts ...Option) error
	Has(ctx context.Context, key string) (bool, error)
	Del(ctx context.Context, key string) error
}

type Depend interface {
	Val(ctx context.Context) interface{}
	Ok(ctx context.Context) error
}

const (
	Redis Type = "redis"
)

type Type string

type Config struct {
	Type
	Addr     string
	Password string
	DB       int
}

var defaultCache ICache

func Cache(ctx context.Context, config *configs.Config) error {
	cfg := &Config{}
	if err := config.Scan("cache", cfg); err != nil {
		return err
	}
	cache, err := NewWithConfig(ctx, cfg)
	if err != nil {
		return err
	}
	defaultCache = cache
	return nil
}

func NewWithConfig(ctx context.Context, cfg *Config, opts ...Option) (ICache, error) {
	redis := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	err := redis.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}
	cache := newRedisCache(redis)
	return cache, nil
}

func Default() ICache {
	return defaultCache
}

type data struct {
	Depend interface{} `json:"depend"`
	Value  interface{} `json:"value"`
}

type StringCache struct {
	String string
}

type IntCache struct {
	Int int
}

type Int64Cache struct {
	Int64 int64
}
