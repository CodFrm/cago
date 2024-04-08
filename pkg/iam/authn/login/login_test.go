package login

import (
	"context"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordLogin(t *testing.T) {
	handler := NewPasswordLogin(func(ctx context.Context, username string) (string, string, error) {
		assert.Equal(t, "user", username)
		b, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		assert.NoError(t, err)
		return "1", string(b), nil
	}, PasswordLoginOptions{}).Login
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("POST", "/", strings.NewReader("username=user&password=123456"))
	c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	uid, err := handler(c)
	assert.NoError(t, err)
	assert.Equal(t, "1", uid)
}

func TestTOTP(t *testing.T) {
	code, err := totp.GenerateCode("qlt6vmy6svfx4bt4rpmisaiyol6hihca", time.Now())
	assert.NoError(t, err)
	assert.NotEmpty(t, code)
	handler := NewTOTP(func(ctx context.Context, username string) (string, string, error) {
		assert.Equal(t, "user", username)
		return "1", "qlt6vmy6svfx4bt4rpmisaiyol6hihca", nil
	}, TOTPOptions{}).Login
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("POST", "/", strings.NewReader("username=user&otpcode="+code))
	c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	uid, err := handler(c)
	assert.NoError(t, err)
	assert.Equal(t, "1", uid)
}

func TestMFA(t *testing.T) {
	errTOPTCodeInvalid := errors.New("code invalid")
	errPasswordNotMatch := errors.New("password not match")
	handler := NewMFA(
		NewPasswordLogin(func(ctx context.Context, username string) (string, string, error) {
			assert.Equal(t, "user", username)
			b, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
			assert.NoError(t, err)
			return "1", string(b), nil
		}, PasswordLoginOptions{ErrPasswordNotMatch: errPasswordNotMatch}),
		NewTOTP(func(ctx context.Context, username string) (string, string, error) {
			assert.Equal(t, "user", username)
			return "1", "qlt6vmy6svfx4bt4rpmisaiyol6hihca", nil
		}, TOTPOptions{ErrTOPTCodeInvalid: errTOPTCodeInvalid}),
	).Login
	code, err := totp.GenerateCode("qlt6vmy6svfx4bt4rpmisaiyol6hihca", time.Now())
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("POST", "/",
		strings.NewReader("username=user&password=123456&otpcode="+code))
	c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	uid, err := handler(c)
	assert.NoError(t, err)
	assert.Equal(t, "1", uid)

	c, _ = gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("POST", "/",
		strings.NewReader("username=user&password=123456&otpcode=123456"))
	c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	uid, err = handler(c)
	assert.Equal(t, errTOPTCodeInvalid, err)
	assert.Empty(t, uid)

	// test password error
	c, _ = gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("POST", "/",
		strings.NewReader("username=user&password=1234567&otpcode="+code))
	c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	uid, err = handler(c)
	assert.Equal(t, errPasswordNotMatch, err)
	assert.Empty(t, uid)
}
