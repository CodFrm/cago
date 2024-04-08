package login

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
)

type TOTPOptions struct {
	ErrUsernameRequired error
	ErrTOPTCodeRequired error
	ErrTOPTCodeInvalid  error
}

// TOTP 动态口令 totp是基于时间的一次性密码算法
type TOTP struct {
	getSecret func(ctx context.Context, username string) (string, string, error)
	options   TOTPOptions
}

func NewTOTP(getSecret func(ctx context.Context, username string) (string, string, error),
	options TOTPOptions) *TOTP {
	return &TOTP{
		getSecret: getSecret,
		options:   options,
	}
}

func (t *TOTP) PreLogin(ctx *gin.Context) error {
	return nil
}

func (t *TOTP) Login(ctx *gin.Context) (string, error) {
	username := ctx.PostForm("username")
	if username == "" {
		return "", t.options.ErrUsernameRequired
	}
	otpcode := ctx.PostForm("otpcode")
	if otpcode == "" {
		return "", t.options.ErrTOPTCodeRequired
	}
	uid, secret, err := t.getSecret(ctx.Request.Context(), username)
	if err != nil {
		return "", err
	}
	if !totp.Validate(otpcode, secret) {
		return "", t.options.ErrTOPTCodeInvalid
	}
	return uid, nil
}
