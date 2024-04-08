package authn

import (
	"context"

	"github.com/codfrm/cago/pkg/iam/sessions"
	"github.com/codfrm/cago/pkg/utils/httputils"
)

var (
	ErrUnauthorized = httputils.NewUnauthorizedError(-1, "未登录")
	ErrTokenInvalid = httputils.NewUnauthorizedError(-2, "token失效")
)

type Options struct {
	sessionManager sessions.HTTPSessionManager
	setContext     func(ctx context.Context, session *sessions.Session) (context.Context, error)
}

type Option func(*Options)

func newOptions(opts ...Option) *Options {
	opt := &Options{}
	for _, o := range opts {
		o(opt)
	}
	return opt
}

func WithSessionManager(session sessions.HTTPSessionManager) Option {
	return func(o *Options) {
		o.sessionManager = session
	}
}

func WithSetContext(setContext func(ctx context.Context, session *sessions.Session) (context.Context, error)) Option {
	return func(o *Options) {
		o.setContext = setContext
	}
}

type LoginOptions struct {
	saveSession SaveSession
}

type LoginOption func(*LoginOptions)

func newLoginOptions(opts ...LoginOption) *LoginOptions {
	opt := &LoginOptions{}
	for _, o := range opts {
		o(opt)
	}
	return opt
}

func WithSaveSession(saveSession SaveSession) LoginOption {
	return func(o *LoginOptions) {
		o.saveSession = saveSession
	}
}
