package gogo

import "context"

// Go 框架处理协程,用于优雅启停
func Go(ctx context.Context, fun func(ctx context.Context) error) error {
	go func() {
		_ = fun(ctx)
	}()
	return nil
}
