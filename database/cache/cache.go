package cache

import (
	"context"
	"github.com/codfrm/cago"
	"github.com/codfrm/cago/configs"
	cache2 "github.com/codfrm/cago/database/cache/cache"
	"github.com/codfrm/cago/database/cache/redis"
	redis2 "github.com/redis/go-redis/v9"
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

var defaultCache *cache

type cache struct {
	cache2.Cache
}

func Cache() cago.Component {
	return &cache{}
}

func (*cache) Start(ctx context.Context, config *configs.Config) error {
	cfg := &Config{}
	if err := config.Scan("cache", cfg); err != nil {
		return err
	}
	c, err := NewWithConfig(ctx, cfg)
	if err != nil {
		return err
	}
	defaultCache = &cache{
		Cache: c,
	}
	return nil
}

func (c *cache) CloseHandle() {
	_ = c.Close()
}

func NewWithConfig(ctx context.Context, cfg *Config, opts ...cache2.Option) (cache2.Cache, error) {
	return redis.NewRedisCache(&redis2.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
}

func Default() cache2.Cache {
	return defaultCache
}
