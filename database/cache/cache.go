package cache

import (
	"context"

	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/database/cache/cache"
	"github.com/codfrm/cago/database/cache/redis"
	redis2 "github.com/go-redis/redis/v8"
)

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

var defaultCache cache.Cache

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

func NewWithConfig(ctx context.Context, cfg *Config, opts ...cache.Option) (cache.Cache, error) {
	return redis.NewRedisCache(&redis2.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
}

func Default() cache.Cache {
	return defaultCache
}
