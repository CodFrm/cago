package source

type Source interface {
	Scan(key string, value interface{}) error
}
