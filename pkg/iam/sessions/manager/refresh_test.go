package manager

import (
	"encoding/json"
	"github.com/codfrm/cago/pkg/iam/authn"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewRefreshHTTPSessionManager(t *testing.T) {
	refreshSession := NewRefreshHTTPSessionManager(NewMemorySessionManager(), NewMemorySessionManager())

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Cookie", "access_token=xxxxx")

	// session不存在
	session, err := refreshSession.GetFromRequest(c)
	assert.Equal(t, authn.ErrSessionNotFound, err)

	// 创建一个新的session
	session, err = refreshSession.Start(c)
	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.Len(t, session.Values, 0)
	assert.Len(t, session.Metadata, 0)
	session.Values["uid"] = 1
	err = refreshSession.SaveToResponse(c, session)
	assert.NoError(t, err)
	// 自动添加了到期时间
	assert.Len(t, session.Values, 1)
	assert.Len(t, session.Metadata, 1)
	assert.Equal(t, 1, session.Values["uid"])
	refreshResp := &RefreshSessionResponse{}
	resp := &httputils.JSONResponse{
		Data: refreshResp,
	}
	err = json.Unmarshal(w.Body.Bytes(), resp)
	assert.NoError(t, err)
	assert.Equal(t, 0, resp.Code)

	// 通过token获取信息
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Cookie", "access_token="+refreshResp.AccessToken)
	session, err = refreshSession.GetFromRequest(c)
	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, 1, session.Values["uid"])

	// 刷新token(不带access_token)
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader("refresh_token="+refreshResp.RefreshToken))
	session, err = refreshSession.GetFromRequest(c)
	assert.Equal(t, authn.ErrSessionNotFound, err)
	assert.Nil(t, session)

	// 刷新token(带access_token)
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader("refresh_token="+refreshResp.RefreshToken))
	c.Request.Header.Set("Cookie", "access_token="+refreshResp.AccessToken)
	c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	session, err = refreshSession.GetFromRequest(c)
	assert.NoError(t, err)
	assert.Equal(t, 1, session.Values["uid"])
	err = refreshSession.SaveToResponse(c, session)
	assert.NoError(t, err)
	refreshResp2 := &RefreshSessionResponse{}
	resp = &httputils.JSONResponse{
		Data: refreshResp2,
	}
	err = json.Unmarshal(w.Body.Bytes(), resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, refreshResp2.AccessToken)
	assert.NotEmpty(t, refreshResp.AccessToken)
	assert.NotEqual(t, refreshResp.AccessToken, refreshResp2.AccessToken)
	assert.NotEqual(t, refreshResp.RefreshToken, refreshResp2.RefreshToken)

	// 使用老的token访问
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	c.Request.Header.Set("Cookie", "access_token="+refreshResp.AccessToken)
	session, err = refreshSession.GetFromRequest(c)
	assert.Equal(t, authn.ErrSessionNotFound, err)
	assert.Nil(t, session)

}
