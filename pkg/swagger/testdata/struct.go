package testdata

type PkgStruct struct {
	Name string `json:"name"`
}

// Nested 嵌套
type Nested struct {
	PkgStruct `json:",inline"`
	Data      string `json:"data"`
}

// TestInfo 测试信息
type TestInfo[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
}
