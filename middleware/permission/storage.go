package permission

import (
	"context"
)

// PolicyStorageFinder 策略存储查找器
type PolicyStorageFinder interface {
	FindPolicyBySubject(ctx context.Context, sub string) ([]*Policy, error)
	FindPolicyByObject(ctx context.Context, sub, obj string) ([]*Policy, error)
}

// PolicyStorage 策略存储
type PolicyStorage interface {
	PolicyStorageFinder
	AddPolicy(ctx context.Context, policy *Policy) error
	RemovePolicy(ctx context.Context, policy *Policy) error
	UpdatePolicy(ctx context.Context, policy *Policy) error
}
