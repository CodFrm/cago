package storage

import (
	"context"
	"github.com/codfrm/cago/middleware/permission"
	"strings"
)

type Group struct {
	permission.PolicyStorage
}

func NewGroup(storage permission.PolicyStorage) *Group {
	return &Group{
		PolicyStorage: storage,
	}
}

func (g *Group) FindPolicyBySubject(ctx context.Context, sub string) ([]*permission.Policy, error) {
	// 先搜索到用户属于哪个组
	userPolicy, err := g.PolicyStorage.FindPolicyBySubject(ctx, sub)
	if err != nil {
		return nil, err
	}
	// 再搜索组的策略
	list := make([]*permission.Policy, 0)
	for _, v := range userPolicy {
		if strings.HasPrefix(v.Object, "group:") {
			policies, err := g.PolicyStorage.FindPolicyBySubject(ctx, v.Object)
			if err != nil {
				return nil, err
			}
			list = append(list, policies...)
		} else {
			list = append(list, v)
		}
	}
	return list, nil
}

func (g *Group) FindPolicyByObject(ctx context.Context, sub, obj string) ([]*permission.Policy, error) {
	// 先搜索到用户属于哪个组
	userPolicy, err := g.PolicyStorage.FindPolicyBySubject(ctx, sub)
	if err != nil {
		return nil, err
	}
	// 再搜索组的策略
	list := make([]*permission.Policy, 0)
	for _, v := range userPolicy {
		if strings.HasPrefix(v.Object, "group:") {
			policies, err := g.PolicyStorage.FindPolicyByObject(ctx, v.Object, obj)
			if err != nil {
				return nil, err
			}
			list = append(list, policies...)
		} else if v.Object == obj {
			list = append(list, v)
		}
	}
	return list, nil
}

// AddUserToGroup 添加用户到组
func (g *Group) AddUserToGroup(ctx context.Context, user, group string) error {
	return g.AddPolicy(ctx, &permission.Policy{
		Subject: user,
		Object:  "group:" + group,
	})
}
