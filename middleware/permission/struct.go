package permission

import (
	"context"
	"errors"
)

var (
	ErrPermissionDenied = errors.New("permission denied")
	ErrPolicyNotFound   = errors.New("not found")
)

// Request 请求
type Request struct {
	Subject     string            // 访问者
	Resource    string            // 资源信息
	Action      string            // 操作
	Environment map[string]string // 环境信息
}

type PolicyEffect string

const (
	Allow PolicyEffect = "allow"
	Deny  PolicyEffect = "deny"
)

// Policy 策略
type Policy struct {
	ID       string       // 策略id
	Subject  string       // 访问者
	Resource string       // 资源信息
	Actions  []string     // 策略的操作类型
	Effect   PolicyEffect // 结果
	Priority int          // 优先级
}

func (p *Policy) InAction(action string) bool {
	for _, v := range p.Actions {
		if v == action {
			return true
		}
	}
	return false
}

// Matcher 策略储存匹配器, 从策略存储中匹配出符合条件的策略
type Matcher interface {
	Match(ctx context.Context, request *Request, policyStorage PolicyStorageFinder) ([]*Policy, error)
}

type FuncMatcher func(ctx context.Context, request *Request, policyStorage PolicyStorageFinder) ([]*Policy, error)

func (f FuncMatcher) Match(ctx context.Context, request *Request, policyStorage PolicyStorageFinder) ([]*Policy, error) {
	return f(ctx, request, policyStorage)
}
