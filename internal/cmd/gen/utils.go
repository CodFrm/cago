package gen

import (
	"go/ast"
	"strings"
)

// 首字母大写
func upperFirstChar(str string) string {
	if len(str) == 0 {
		return str
	}
	return strings.ToUpper(str[0:1]) + str[1:]
}

// 首字母小写
func lowerFirstChar(str string) string {
	if len(str) == 0 {
		return str
	}
	return strings.ToLower(str[:1]) + str[1:]
}

func getComment(decl *ast.GenDecl, typeSpec *ast.TypeSpec) string {
	comment := ""
	if decl.Doc != nil {
		comment = decl.Doc.Text()
		comment = strings.TrimSpace(strings.TrimPrefix(comment, typeSpec.Name.Name))
	}
	return comment
}
