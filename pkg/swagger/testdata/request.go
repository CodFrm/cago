package testdata

import (
	"github.com/codfrm/cago/pkg/swagger/testdata/pkg"
	"github.com/codfrm/cago/server/mux"
)

// TestInfo 测试信息
type TestInfo struct {
	Name string `json:"name"`
	Pkg  pkg.Enum
	// Nested 嵌套
	Nested Nested `json:"nested"`
}

// TestRequest test
type TestRequest struct {
	mux.Meta `path:"/test" method:"GET"`
	Name     string `json:"name"` // 名字
	// Age 年龄
	Age  int       `json:"age"`
	Enum *pkg.Enum `json:"enum"` // 嵌套类型
}

type TestResponse struct {
	List      []interface{}          `json:"list"`
	Map       map[string]interface{} `json:"map"`
	Interface interface{}            `json:"interface"`
	Inline    struct {
		Name string `json:"name"`
	}
	// Info 123
	Info *TestInfo `json:"info"`
	// PkgStruct 123
	PkgStruct *PkgStruct
	Enum      *pkg.Enum
}
