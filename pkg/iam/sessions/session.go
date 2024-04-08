package sessions

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session expired")
)

type Session struct {
	// ID 会话ID
	ID string
	// Metadata 会话元数据
	Metadata map[string]interface{} `json:"metadata"`
	// Values 会话值
	Values map[string]interface{} `json:"values"`
}

// SessionManager 会话管理
type SessionManager interface {
	// Start 开始会话 返回一个新的会话
	Start(ctx context.Context) (*Session, error)
	// Get 获取会话 如果会话不存在则返回ErrSessionNotFound
	Get(ctx context.Context, id string) (*Session, error)
	// Save 保存会话
	Save(ctx context.Context, session *Session) error
	// Delete 删除会话
	Delete(ctx context.Context, id string) error
	// Refresh 刷新会话
	Refresh(ctx context.Context, session *Session) error
}

type HTTPSessionManager interface {
	SessionManager
	// GetFromRequest 从请求中获取会话 如果会话不存在则返回ErrSessionNotFound
	GetFromRequest(ctx *gin.Context) (*Session, error)
	// SaveToResponse 保存session到响应
	SaveToResponse(ctx *gin.Context, session *Session) error
}
