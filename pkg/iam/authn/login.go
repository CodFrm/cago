package authn

import (
	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/cago/pkg/iam/sessions"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/gin-gonic/gin"
)

// LoginHandler 登录操作 登录成功后返回用户信息
type LoginHandler func(ctx *gin.Context) (*User, error)

// SaveSession 保存session
type SaveSession func(ctx *gin.Context, user *User, session *sessions.Session) error

var (
	UsernameNotFound           = 10000
	UsernameOrPasswordRequired = 10001
	PasswordWrong              = 10002
)

func (a *Authn) LoginByPassword() gin.HandlerFunc {
	return a.LoginBy(func(ctx *gin.Context) (*User, error) {
		username := ctx.PostForm("username")
		password := ctx.PostForm("password")
		if username == "" || password == "" {
			return nil, i18n.NewError(ctx, UsernameOrPasswordRequired)
		}
		user, err := a.database.GetUserByUsername(ctx, username, WithPassword())
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, i18n.NewError(ctx, UsernameNotFound)
		}
		if err := user.CheckPassword(password); err != nil {
			return nil, i18n.NewError(ctx, PasswordWrong)
		}
		return user, nil
	})
}

func (a *Authn) LoginBy(login LoginHandler, opts ...LoginOption) gin.HandlerFunc {
	options := newLoginOptions(opts...)
	return func(ctx *gin.Context) {
		httputils.Handle(ctx, func() interface{} {
			user, err := login(ctx)
			if err != nil {
				return err
			}
			// 登录成功，设置session
			session, err := a.options.sessionManager.Start(ctx)
			if err != nil {
				return err
			}
			session.Values["uid"] = user.ID
			if options.saveSession != nil {
				if err := options.saveSession(ctx, user, session); err != nil {
					return err
				}
			}
			return a.options.sessionManager.SaveToResponse(ctx, session)
		})
	}
}
