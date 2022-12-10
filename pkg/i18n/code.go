package i18n

import (
	"context"
	"fmt"
	"net/http"

	"github.com/codfrm/cago/pkg/utils/httputils"
)

var DefaultLang = "zh-cn"

var langs = map[string]map[int]string{}

func Register(lang string, code map[int]string) {
	if _, ok := langs[lang]; ok {
		// append
		for k, v := range code {
			langs[lang][k] = v
		}
	} else {
		langs[lang] = code
	}
}

// NewError 参数校验错误
func NewError(ctx context.Context, code int, v ...interface{}) error {
	return httputils.NewError(http.StatusBadRequest, code, Printf(ctx, code, v...))
}

// NewInternalError 构造内部错误
func NewInternalError(ctx context.Context, code int, v ...interface{}) error {
	return httputils.NewError(http.StatusInternalServerError, code, Printf(ctx, code, v...))
}

func NewErrorWithStatus(ctx context.Context, status int, code int, v ...interface{}) error {
	return httputils.NewError(status, code, Printf(ctx, code, v...))
}

func Printf(ctx context.Context, code int, v ...interface{}) string {
	if _, ok := langs[DefaultLang]; !ok {
		return fmt.Sprintf("code %d not found", code)
	}
	return fmt.Sprintf(langs[DefaultLang][code], v...)
}
