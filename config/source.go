package config

type Source interface {
	Scan(key string, value interface{}) error
}
