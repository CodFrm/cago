package permission

import (
	"context"
)

type equalMatcher struct {
}

func NewEqualMatcher() Matcher {
	return &equalMatcher{}
}

func (m *equalMatcher) Match(ctx context.Context, request *Request, policyStorage PolicyStorageFinder) ([]*Policy, error) {
	// 查询出相关策略
	policies, err := policyStorage.FindPolicyByObject(ctx, request.Subject, request.Object)
	if err != nil {
		return nil, err
	}
	matchedPolicies := make([]*Policy, 0)
	for _, policy := range policies {
		if policy.Object == request.Object {
			if policy.InAction(request.Action) {
				matchedPolicies = append(matchedPolicies, policy)
			}
		}
	}
	return matchedPolicies, nil
}

type orMatcher struct {
	matches []Matcher
}

func NewOrMatcher(matches ...Matcher) Matcher {
	return &orMatcher{
		matches: matches,
	}
}

func (m *orMatcher) Match(ctx context.Context, request *Request, policyStorage PolicyStorageFinder) ([]*Policy, error) {
	for _, v := range m.matches {
		if policies, err := v.Match(ctx, request, policyStorage); err != nil {
			return nil, err
		} else if len(policies) > 0 {
			return policies, nil
		}
	}
	return nil, nil
}
