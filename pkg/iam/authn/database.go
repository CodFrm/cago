package authn

import (
	"context"
)

type GetUserOptions struct {
	// 带上密码
	WithPassword bool
	// Metadata 用户元数据
	Metadata map[UserMetadata]interface{} `json:"metadata,omitempty"`
}

type GetUserOption func(*GetUserOptions)

func NewGetUserOptions(opts ...GetUserOption) *GetUserOptions {
	o := &GetUserOptions{}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func WithPassword() GetUserOption {
	return func(o *GetUserOptions) {
		o.WithPassword = true
	}
}

type Database interface {
	// Register 注册用户
	Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
	// GetUserByUsername 通过用户名获取用户
	GetUserByUsername(ctx context.Context, username string, opts ...GetUserOption) (*User, error)
	// GetUserByID 通过用户ID获取用户
	GetUserByID(ctx context.Context, userID string, opts ...GetUserOption) (*User, error)
	// GetUserByWhere 通过条件获取用户
	GetUserByWhere(ctx context.Context, where map[string]interface{}, opts ...GetUserOption) (*User, error)
	// UpdateUser 更新用户信息
	//UpdateUser(ctx context.Context, userID string, user *UpdateUserRequest) error
}
