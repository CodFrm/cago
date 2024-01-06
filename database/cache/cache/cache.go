package cache

import (
	"context"
	"errors"
)

var (
	ErrNil = errors.New("cache: nil")
)

type Cache interface {
	GetOrSet(ctx context.Context, key string, set func() (interface{}, error), opts ...Option) Value
	Set(ctx context.Context, key string, val interface{}, opts ...Option) Value
	Get(ctx context.Context, key string, opts ...Option) Value
	Has(ctx context.Context, key string) (bool, error)
	Del(ctx context.Context, key string) error
	Close() error
}

type Depend interface {
	Val(ctx context.Context) interface{}
	Valid(ctx context.Context) error
}

type Value interface {
	Result() (string, error)
	Err() error
	Scan(v interface{}) error
	Bytes() ([]byte, error)
	Int64() (int64, error)
	Bool() (bool, error)
}

func IsNil(err error) bool {
	return errors.Is(err, ErrNil)
}
