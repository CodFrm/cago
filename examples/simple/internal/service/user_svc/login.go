package user_svc

import (
	"context"
	"github.com/codfrm/cago/examples/simple/internal/model"
	"github.com/codfrm/cago/examples/simple/internal/model/entity/user_entity"
	"github.com/codfrm/cago/examples/simple/internal/pkg/code"
	"github.com/codfrm/cago/examples/simple/internal/repository/user_repo"
	"github.com/codfrm/cago/middleware/sessions"
	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/pkg/opentelemetry/trace"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"time"

	api "github.com/codfrm/cago/examples/simple/internal/api/user"
)

type LoginSvc interface {
	// Register 注册
	Register(ctx context.Context, req *api.RegisterRequest) (*api.RegisterResponse, error)
	// Login 登录
	Login(ctx *gin.Context, req *api.LoginRequest) (*api.LoginResponse, error)
	// Logout 登出
	Logout(ctx *gin.Context, req *api.LogoutRequest) (*api.LogoutResponse, error)
	// Middleware 鉴权中间件
	Middleware() gin.HandlerFunc
	// Get 从上下文中获取用户登录信息
	Get(ctx context.Context) *model.AuthInfo
}

type loginSvc struct {
}

var defaultLogin = &loginSvc{}

func Login() LoginSvc {
	return defaultLogin
}

// Register 注册
func (l *loginSvc) Register(ctx context.Context, req *api.RegisterRequest) (*api.RegisterResponse, error) {
	// 查找相同用户名
	user, err := user_repo.User().FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return nil, i18n.NewError(ctx, code.UsernameAlreadyExists)
	}
	// 创建用户
	user = &user_entity.User{
		Username:   req.Username,
		Status:     consts.ACTIVE,
		Createtime: time.Now().Unix(),
	}
	if err := user_repo.User().Create(ctx, user); err != nil {
		return nil, err
	}
	return &api.RegisterResponse{}, nil
}

// Login 登录
func (l *loginSvc) Login(ctx *gin.Context, req *api.LoginRequest) (*api.LoginResponse, error) {
	// 查找用户
	user, err := user_repo.User().FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if err := user.Check(ctx); err != nil {
		return nil, err
	}
	// 设置session
	sessions.Ctx(ctx).Set("user_id", user.ID)
	if err := sessions.Ctx(ctx).Save(); err != nil {
		return nil, err
	}
	return &api.LoginResponse{
		Username: user.Username,
	}, nil
}

// Logout 登出
func (l *loginSvc) Logout(ctx *gin.Context, req *api.LogoutRequest) (*api.LogoutResponse, error) {
	sessions.Ctx(ctx).Delete("user_id")
	if err := sessions.Ctx(ctx).Save(); err != nil {
		return nil, err
	}
	return &api.LogoutResponse{}, nil
}

func (l *loginSvc) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, ok := sessions.Ctx(c).Get("user_id").(int64)
		if !ok {
			httputils.HandleResp(c, i18n.NewError(c.Request.Context(), code.UserNotLogin))
			return
		}
		ctx, err := l.SetUser(c.Request.Context(), uid)
		if err != nil {
			httputils.HandleResp(c, err)
			return
		}
		c.Request = c.Request.WithContext(ctx)
	}
}

func (l *loginSvc) SetUser(ctx context.Context, uid int64) (context.Context, error) {
	user, err := user_repo.User().Find(ctx, uid)
	if err != nil {
		return nil, err
	}
	if err := user.Check(ctx); err != nil {
		return nil, err
	}
	// 设置用户信息,链路追踪和日志也添加上用户信息
	authInfo := &model.AuthInfo{
		UserID:   uid,
		Username: user.Username,
	}
	trace.SpanFromContext(ctx).SetAttributes(
		attribute.Int64("user_id", user.ID),
	)
	return context.WithValue(
		logger.ContextWithLogger(ctx, logger.Ctx(ctx).
			With(zap.Int64("user_id", user.ID))),
		model.AuthInfo{}, authInfo), nil
}

func (l *loginSvc) Get(ctx context.Context) *model.AuthInfo {
	val := ctx.Value(model.AuthInfo{})
	if val == nil {
		return nil
	}
	return val.(*model.AuthInfo)
}
