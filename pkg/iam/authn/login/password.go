package login

import (
	"context"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type PasswordLoginOptions struct {
	ErrUsernameRequired error
	ErrPasswordRequired error
	ErrPasswordNotMatch error
}

// PasswordLogin 密码登录
type PasswordLogin struct {
	getPassword func(ctx context.Context, username string) (string, string, error)
	options     PasswordLoginOptions
}

func NewPasswordLogin(getPassword func(ctx context.Context, username string) (string, string, error),
	options PasswordLoginOptions) *PasswordLogin {
	return &PasswordLogin{
		getPassword: getPassword,
		options:     options,
	}
}

func (p *PasswordLogin) PreLogin(ctx *gin.Context) error {
	return nil
}

func (p *PasswordLogin) Login(ctx *gin.Context) (string, error) {
	username := ctx.PostForm("username")
	if username == "" {
		return "", p.options.ErrUsernameRequired
	}
	password := ctx.PostForm("password")
	if password == "" {
		return "", p.options.ErrPasswordRequired
	}
	// 获取用户密码
	id, pwd, err := p.getPassword(ctx.Request.Context(), username)
	if err != nil {
		return "", err
	}
	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(pwd), []byte(password)); err != nil {
		return "", p.options.ErrPasswordNotMatch
	}
	return id, nil
}
