package sync

import (
	"errors"
	"time"
)

var (
	ErrTryLockTimeout = errors.New("尝试获取锁超时")
	ErrLockOccurred   = errors.New("锁已经被占用")
	ErrLockNotExists  = errors.New("锁不存在")
)

type LockOption func(*LockOptions)

type LockOptions struct {
	// 设置获取锁的堵塞超时时间
	timeout time.Duration
}

// WithLockTimeout 设置超时时间
func WithLockTimeout(timeout time.Duration) LockOption {
	return func(lo *LockOptions) {
		lo.timeout = timeout
	}
}
