package token_auth

import (
	"context"
	"strings"

	"github.com/codfrm/cago/pkg/sync"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/gin-gonic/gin"
)

var (
	ErrUnauthorized        = httputils.NewUnauthorizedError(-1, "未登录")
	ErrTokenInvalid        = httputils.NewUnauthorizedError(-2, "token失效")
	ErrTokenExpired        = httputils.NewUnauthorizedError(-3, "token过期")
	ErrRefreshTokenInvalid = httputils.NewUnauthorizedError(-4, "refresh_token失效")
	ErrRefreshTokenExpired = httputils.NewUnauthorizedError(-5, "refresh_token过期")
)

type Options struct {
	getAccessToken func(ctx *gin.Context) (string, error)
	handlerError   func(ctx *gin.Context, err error) error
	storage        Storage
	lock           sync.Locker
	setContext     func(ctx context.Context, accessToken *AccessToken) (context.Context, error)
	tokenExpired   int64
	refreshExpired int64
}

type Option func(*Options)

func newOptions(opts ...Option) *Options {
	opt := &Options{
		getAccessToken: defaultGetAccessToken,
		handlerError: func(ctx *gin.Context, err error) error {
			return err
		},
		lock:           sync.NewLocker("token_auth"),
		tokenExpired:   7200,
		refreshExpired: 3600 * 24 * 30,
	}
	for _, o := range opts {
		o(opt)
	}
	return opt
}

func WithGetAccessToken(getToken func(ctx *gin.Context) (string, error)) Option {
	return func(o *Options) {
		o.getAccessToken = getToken
	}
}

func WithHandlerError(handlerError func(ctx *gin.Context, err error) error) Option {
	return func(o *Options) {
		o.handlerError = handlerError
	}
}

func WithStorage(storage Storage) Option {
	return func(o *Options) {
		o.storage = storage
	}
}

func WithLock(lock sync.Locker) Option {
	return func(o *Options) {
		o.lock = lock
	}
}

func WithSetContext(setContext func(ctx context.Context, accessToken *AccessToken) (context.Context, error)) Option {
	return func(o *Options) {
		o.setContext = setContext
	}
}

func WithTokenExpired(tokenExpired int64, refreshExpired int64) Option {
	return func(o *Options) {
		o.tokenExpired = tokenExpired
		o.refreshExpired = refreshExpired
	}
}

func defaultGetAccessToken(ctx *gin.Context) (string, error) {
	// 找到access_token
	accessToken, _ := ctx.Cookie("access_token")
	if accessToken == "" {
		accessToken = ctx.GetHeader("Authorization")
		if accessToken != "" {
			accessToken = strings.TrimPrefix(accessToken, "Bearer ")
		}
	}
	return accessToken, nil
}
