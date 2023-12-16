package permission

import (
	"context"
	"sort"
)

type Permission struct {
	policyStorage PolicyStorage
	options       *Options
}

func NewPermission(policyStorage PolicyStorage, opts ...Option) *Permission {
	options := &Options{
		storageMatcher: &equalMatcher{},
	}
	for _, opt := range opts {
		opt(options)
	}
	return &Permission{
		policyStorage: policyStorage,
		options:       options,
	}
}

// Check 检查权限
func (p *Permission) Check(ctx context.Context, request *Request, opts ...CheckOption) error {
	options := &CheckOptions{}
	for _, v := range opts {
		v(options)
	}
	policies, err := p.queryPolicies(ctx, request, options)
	if err != nil {
		return err
	}
	if len(policies) == 0 {
		return ErrPermissionDenied
	} else if len(policies) == 1 {
		if policies[0].Effect == Deny {
			return ErrPermissionDenied
		}
	} else {
		// 按优先级排序
		sort.SliceStable(policies, func(i, j int) bool {
			return policies[i].Priority < policies[j].Priority
		})
		if policies[0].Effect == Deny {
			return ErrPermissionDenied
		}
	}
	return nil
}

// QueryPolicies 查询出符合条件的策略
func (p *Permission) QueryPolicies(ctx context.Context, request *Request, opts ...CheckOption) ([]*Policy, error) {
	options := &CheckOptions{}
	for _, v := range opts {
		v(options)
	}
	return p.queryPolicies(ctx, request, options)
}

func (p *Permission) queryPolicies(ctx context.Context, request *Request, options *CheckOptions) ([]*Policy, error) {
	// 查询出相关策略
	storageMatcher := p.options.storageMatcher
	if options.storageMatcher != nil {
		storageMatcher = options.storageMatcher
	}
	policies, err := storageMatcher.Match(ctx, request, p.policyStorage)
	if err != nil {
		return nil, err
	}
	return policies, nil
}
