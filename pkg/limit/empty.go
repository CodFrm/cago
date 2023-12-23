package limit

import "context"

type Empty struct {
}

func (e *Empty) Take(ctx context.Context, key string) (func() error, error) {
	return func() error {
		return nil
	}, nil
}

func (e *Empty) FuncTake(ctx context.Context, key string, f func() (interface{}, error)) (interface{}, error) {
	return f()
}

// NewEmpty 创建空限流器
func NewEmpty() Limit {
	return &Empty{}
}
