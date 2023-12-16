package permission_test

import (
	"context"
	"github.com/codfrm/cago/middleware/permission"
	"github.com/codfrm/cago/middleware/permission/memory"
	"github.com/codfrm/cago/middleware/permission/storage"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestACL(t *testing.T) {
	ctx := context.Background()
	storage := memory.NewMemory()
	p := permission.NewPermission(storage)
	_ = storage.AddPolicy(ctx, &permission.Policy{
		Subject: "user1",
		Object:  "/file",
		Actions: []string{"read"},
		Effect:  permission.Allow,
	})
	err := p.Check(ctx, &permission.Request{
		Subject: "user1",
		Object:  "/file",
		Action:  "read",
	})
	assert.Nil(t, err)
	err = p.Check(ctx, &permission.Request{
		Subject: "user1",
		Object:  "/file",
		Action:  "write",
	})
	assert.Equal(t, err, permission.ErrPermissionDenied)
	// 带超级管理员的匹配器
	_ = storage.AddPolicy(ctx, &permission.Policy{
		Subject: "admin",
		Effect:  permission.Allow,
	})
	_ = storage.AddPolicy(ctx, &permission.Policy{
		Subject:  "admin",
		Object:   "/file2",
		Actions:  []string{"write"},
		Effect:   permission.Deny,
		Priority: -1,
	})
	p = permission.NewPermission(storage, permission.WithStorageMatcher(
		permission.NewOrMatcher(
			permission.FuncMatcher(func(ctx context.Context, request *permission.Request, policyStorage permission.PolicyStorageFinder) ([]*permission.Policy, error) {
				policies, err := policyStorage.FindPolicyBySubject(ctx, request.Subject)
				if err != nil {
					return nil, err
				}
				adminPolicies := make([]*permission.Policy, 0)
				for _, policy := range policies {
					if policy.Subject == "admin" &&
						(policy.Object == "" || policy.Object == request.Object) {
						if len(policy.Actions) == 0 || policy.InAction(request.Action) {
							adminPolicies = append(adminPolicies, policy)
						}
					}
				}
				return adminPolicies, nil
			}),
			permission.NewEqualMatcher()),
	))
	err = p.Check(ctx, &permission.Request{
		Subject: "admin",
		Object:  "/file",
		Action:  "write",
	})
	assert.Nil(t, err)
	err = p.Check(ctx, &permission.Request{
		Subject: "admin",
		Object:  "/file2",
		Action:  "write",
	})
	assert.Equal(t, err, permission.ErrPermissionDenied)
	err = p.Check(ctx, &permission.Request{
		Subject: "admin",
		Object:  "/file2",
		Action:  "read",
	})
	assert.Nil(t, err)
}

func TestRBAC(t *testing.T) {
	ctx := context.Background()
	memStorage := memory.NewMemory()
	group := storage.NewGroup(memStorage)
	p := permission.NewPermission(group)
	_ = group.AddPolicy(ctx, &permission.Policy{
		Subject: "group:group1",
		Object:  "/file",
		Actions: []string{"read"},
	})
	_ = group.AddPolicy(ctx, &permission.Policy{
		Subject: "group:group2",
		Object:  "/file",
		Actions: []string{"write"},
	})
	_ = group.AddUserToGroup(ctx, "user1", "group1")
	err := p.Check(ctx, &permission.Request{
		Subject: "user1",
		Object:  "/file",
		Action:  "read",
	})
	assert.Nil(t, err)
	err = p.Check(ctx, &permission.Request{
		Subject: "user1",
		Object:  "/file",
		Action:  "write",
	})
	assert.Equal(t, err, permission.ErrPermissionDenied)

}
