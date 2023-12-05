package source

import "errors"

var (
	ErrNotFound = errors.New("config key not found")
)

type Source interface {
	Scan(key string, value interface{}) error
	Has(key string) (bool, error)
}
