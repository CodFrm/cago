package redis

import "github.com/go-redis/redis/v8"

func Nil(err error) bool {
	return err == redis.Nil
}
