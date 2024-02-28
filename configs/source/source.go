package source

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("config key not found")
)

type Event int

const (
	Update Event = iota + 1
	Delete
)

type Source interface {
	Scan(ctx context.Context, key string, value interface{}) error
	Has(ctx context.Context, key string) (bool, error)
	Watch(ctx context.Context, key string, callback func(event Event)) error
}
