package redis

import (
	"context"

	"github.com/codfrm/cago/configs"
	"github.com/go-redis/redis/v8"
)

var defaultRedis *redis.Client

type Config struct {
	Addr     string
	Password string
	DB       int
}

func Redis(ctx context.Context, config *configs.Config) error {
	cfg := &Config{}
	if err := config.Scan("redis", cfg); err != nil {
		return err
	}
	ret := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	err := ret.Ping(context.Background()).Err()
	if err != nil {
		return err
	}
	defaultRedis = ret
	return nil
}

func Default() *redis.Client {
	return defaultRedis
}

func Ctx(ctx context.Context) CtxRedis {
	return &ctxRedis{
		client: defaultRedis,
		ctx:    ctx,
	}
}
