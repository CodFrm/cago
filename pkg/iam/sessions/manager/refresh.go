package manager

import (
	"errors"
	"github.com/codfrm/cago/pkg/iam/sessions"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)

type RefreshHTTPSessionManagerOptions struct {
	// token有效期
	TokenDuration int
	// 刷新有效期
	RefreshDuration int
	// 必须包含access_token 刷新token时必须要有access_token
	RequireAccessToken bool
}

type RefreshHTTPSessionManagerOption func(*RefreshHTTPSessionManagerOptions)

// refreshHTTPSessionManager 刷新session管理器
type refreshHTTPSessionManager struct {
	sessions.SessionManager
	originSessionManager  sessions.SessionManager
	refreshSessionManager sessions.SessionManager
	options               *RefreshHTTPSessionManagerOptions
}

// NewRefreshHTTPSessionManager 创建一个拥有两个token的session管理器
// 拥有一个access_token和一个refresh_token 用于token双刷
func NewRefreshHTTPSessionManager(sessionManager, refreshSessionManager sessions.SessionManager,
	opts ...RefreshHTTPSessionManagerOption) sessions.HTTPSessionManager {
	options := &RefreshHTTPSessionManagerOptions{
		TokenDuration:   7200,
		RefreshDuration: 86400 * 30,
	}
	for _, opt := range opts {
		opt(options)
	}
	return &refreshHTTPSessionManager{
		SessionManager:        NewExpireSessionManager(options.TokenDuration, sessionManager),
		originSessionManager:  sessionManager,
		refreshSessionManager: NewExpireSessionManager(options.RefreshDuration, refreshSessionManager),
		options:               options,
	}
}

func (h *refreshHTTPSessionManager) GetFromRequest(ctx *gin.Context) (*sessions.Session, error) {
	id, err := ctx.Cookie("access_token")
	if err != nil {
		return nil, sessions.ErrSessionNotFound
	}
	if id == "" {
		if h.options.RequireAccessToken {
			return nil, sessions.ErrSessionNotFound
		}
		return h.Start(ctx)
	}
	session, err := h.SessionManager.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return session, nil
}

type RefreshSessionResponse struct {
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
	Expire        int64  `json:"expire"`
	RefreshExpire int64  `json:"refresh_expire"`
}

func (h *refreshHTTPSessionManager) SaveToResponse(ctx *gin.Context, session *sessions.Session) error {
	if session.ID == "" {
		// 如果是新的session，创建一个refresh session与之绑定
		refreshSession, err := h.refreshSessionManager.Start(ctx)
		if err != nil {
			return err
		}
		if err := h.SessionManager.Save(ctx, session); err != nil {
			return err
		}
		refreshSession.Values["access_token"] = session.ID
		if err := h.refreshSessionManager.Save(ctx, refreshSession); err != nil {
			return err
		}
		// 返回
		ctx.JSON(http.StatusOK, httputils.JSONResponse{
			Code: 0,
			Data: &RefreshSessionResponse{
				AccessToken:   session.ID,
				RefreshToken:  refreshSession.ID,
				Expire:        session.Metadata["expire"].(int64),
				RefreshExpire: refreshSession.Metadata["expire"].(int64),
			},
		})
		return nil
	}
	// 判断有没有刷新token，如果有则做一次刷新
	refreshID := ctx.PostForm("refresh_token")
	if refreshID != "" {
		refreshSession, err := h.refreshSessionManager.Get(ctx, refreshID)
		if err != nil {
			return err
		}
		// 刷新token
		if err := h.refreshSessionManager.Refresh(ctx, refreshSession); err != nil {
			return err
		}
		// 刷新access token
		if err := h.SessionManager.Refresh(ctx, session); err != nil {
			return err
		}
		// 返回
		ctx.JSON(http.StatusOK, httputils.JSONResponse{
			Code: 0,
			Data: &RefreshSessionResponse{
				AccessToken:   session.ID,
				RefreshToken:  refreshSession.ID,
				Expire:        session.Metadata["expire"].(int64),
				RefreshExpire: refreshSession.Metadata["expire"].(int64),
			},
		})
		return nil
	}
	return ErrRefreshTokenNotFound
}
