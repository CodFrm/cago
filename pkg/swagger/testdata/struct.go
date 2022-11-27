package testdata

type PkgStruct struct {
	Name string `json:"name"`
}

// Nested 嵌套
type Nested struct {
	PkgStruct `json:",inline"`
	Data      string `json:"data"`
}
