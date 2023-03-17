package redis

import "github.com/redis/go-redis/v9"

func Nil(err error) bool {
	return err == redis.Nil
}
