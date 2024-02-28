package muxtest

import (
	"net/http"
	"net/http/httptest"

	"github.com/codfrm/cago/pkg/utils/validator"
	"github.com/codfrm/cago/server/mux"
	"github.com/codfrm/cago/server/mux/muxclient"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type TestMux struct {
	*mux.Router
	*muxclient.Client
}

type testTransport struct {
	r *gin.Engine
}

func (t *testTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	t.r.ServeHTTP(w, r)
	return w.Result(), nil
}

func NewTestMux() *TestMux {
	r := gin.Default()
	var err error
	binding.Validator, err = validator.NewValidator()
	if err != nil {
		panic(err)
	}
	// ginContext支持fallback
	r.ContextWithFallback = true
	return &TestMux{
		Router: &mux.Router{
			Routes:  &mux.Routes{IRoutes: r},
			IRouter: r,
		},
		Client: muxclient.NewClient("", muxclient.WithClient(&http.Client{
			Transport: &testTransport{r: r},
		})),
	}
}
