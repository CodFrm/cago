package cache

import (
	"github.com/codfrm/cago/database/cache/cache"
)

var (
	Expiration = cache.Expiration
	WithDepend = cache.WithDepend
)

func IsNil(err error) bool {
	return cache.IsNil(err)
}
