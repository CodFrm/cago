package testdata

import (
	"github.com/codfrm/cago/server/mux"
)

// TestRequest test
type TestRequest struct {
	mux.Meta `path:"/test/:uid" method:"GET"`
}
type Item struct {
	ID int64
}

type TestResponse struct {
	TestInfo[*Item] `json:",inline"`
}
