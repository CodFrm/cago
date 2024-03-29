package limit

import "context"

// Limit 限流器接口
type Limit interface {
	Take(ctx context.Context, key string) (func() error, error)
	FuncTake(ctx context.Context, key string, f func() (interface{}, error)) (interface{}, error)
}

// CombinationLimit 组合限流器 可以组合多个限流器
type CombinationLimit struct {
	limits []Limit
}

// NewCombinationLimit 创建组合限流器
func NewCombinationLimit(limits ...Limit) *CombinationLimit {
	return &CombinationLimit{
		limits: limits,
	}
}

func (c *CombinationLimit) cancels(cancels []func() error) func() error {
	return func() error {
		var lastErr error
		for _, cancel := range cancels {
			if err := cancel(); err != nil {
				lastErr = err
			}
		}
		return lastErr
	}
}

func (c *CombinationLimit) Take(ctx context.Context, key string) (func() error, error) {
	cancels := make([]func() error, 0)
	for _, limit := range c.limits {
		f, err := limit.Take(ctx, key)
		if err != nil {
			return c.cancels(cancels), err
		} else if f != nil {
			cancels = append(cancels, f)
		}
	}
	return c.cancels(cancels), nil
}

func (c *CombinationLimit) FuncTake(ctx context.Context, key string, f func() (interface{}, error)) (interface{}, error) {
	cancel, err := c.Take(ctx, key)
	if err != nil {
		if err := cancel(); err != nil {
			return nil, err
		}
		return nil, err
	}
	resp, err := f()
	if err != nil {
		if err := cancel(); err != nil {
			return nil, err
		}
		return nil, err
	}
	return resp, nil
}
