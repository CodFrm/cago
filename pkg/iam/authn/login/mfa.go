package login

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	PreLogin(ctx *gin.Context) error
	Login(ctx *gin.Context) (string, error)
}

var ErrMFAUserIDNotMatch = errors.New("mfa user id not match")

type MFA struct {
	factors []Handler
}

// NewMFA 多因素认证
func NewMFA(factors ...Handler) *MFA {
	return &MFA{
		factors: factors,
	}
}

func (m *MFA) PreLogin(ctx *gin.Context) error {
	for _, f := range m.factors {
		if err := f.PreLogin(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (m *MFA) Login(ctx *gin.Context) (string, error) {
	lastUid := ""
	for _, f := range m.factors {
		uid, err := f.Login(ctx)
		if err != nil {
			return "", err
		}
		if lastUid != "" && lastUid != uid {
			return "", ErrMFAUserIDNotMatch
		}
		lastUid = uid
	}
	return lastUid, nil
}
