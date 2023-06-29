package token_auth

import (
	"context"
	"time"

	"github.com/codfrm/cago/pkg/utils"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/gin-gonic/gin"
)

var (
	accessTokenKey = "accessTokenKey"
)

type TokenAuth struct {
	options *Options
}

func NewTokenAuth(option ...Option) *TokenAuth {
	opts := newOptions(option...)
	return &TokenAuth{
		options: opts,
	}
}

func (t *TokenAuth) Middleware(force bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, err := t.options.getAccessToken(ctx)
		if err != nil {
			httputils.HandleResp(ctx, err)
			return
		}
		if accessToken == "" {
			httputils.HandleResp(ctx, t.options.handlerError(ctx, ErrUnauthorized))
			return
		}
		// 验证token
		m, err := t.options.storage.FindByAccessToken(ctx, accessToken)
		if err != nil {
			httputils.HandleResp(ctx, err)
			return
		}
		if m == nil {
			if force {
				httputils.HandleResp(ctx, t.options.handlerError(ctx, ErrTokenInvalid))
			}
			return
		}
		if m.ExpireAt < time.Now().Unix() {
			if force {
				httputils.HandleResp(ctx, t.options.handlerError(ctx, ErrTokenExpired))
			}
			return
		}
		reqCtx := ctx.Request.Context()
		if t.options.setContext != nil {
			reqCtx, err = t.options.setContext(reqCtx, m)
			if err != nil {
				httputils.HandleResp(ctx, err)
				return
			}
		}
		ctx.Request = ctx.Request.WithContext(context.WithValue(reqCtx, accessTokenKey, m))
	}
}

func (t *TokenAuth) Generate(ctx context.Context) *AccessToken {
	return &AccessToken{
		AccessToken:  utils.RandString(32, utils.Mix),
		RefreshToken: utils.RandString(32, utils.Mix),
		KvMap:        map[string]string{},
		ExpireAt:     time.Now().Unix() + t.options.tokenExpired,
		RefreshAt:    time.Now().Unix() + t.options.refreshExpired,
	}
}

func (t *TokenAuth) Save(ctx context.Context, accessToken *AccessToken) error {
	return t.options.storage.Save(ctx, accessToken)
}

func (t *TokenAuth) Delete(ctx context.Context, accessToken *AccessToken) error {
	return t.options.storage.Delete(ctx, accessToken)
}

func (t *TokenAuth) Get(ctx context.Context) *AccessToken {
	m, ok := ctx.Value(accessTokenKey).(*AccessToken)
	if !ok {
		return nil
	}
	return m
}

func (t *TokenAuth) Refresh(ctx context.Context, refreshToken string) (*AccessToken, error) {
	if err := t.options.lock.LockKey(refreshToken); err != nil {
		return nil, err
	}
	defer func() {
		_ = t.options.lock.UnlockKey(refreshToken)
	}()
	m, err := t.options.storage.FindByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, ErrRefreshTokenInvalid
	}
	if m.RefreshAt < time.Now().Unix() {
		return nil, ErrRefreshTokenExpired
	}
	// 删除原来的
	if err := t.options.storage.Delete(ctx, m); err != nil {
		return nil, err
	}
	// 生成新的
	generate := t.Generate(ctx)
	generate.KvMap = m.KvMap
	return m, t.options.storage.Save(ctx, m)
}
