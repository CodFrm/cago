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

// 下划线转驼峰
func toCamel(str string) string {
	if str == "id" {
		return "ID"
	}
	var result string
	for _, v := range strings.Split(str, "_") {
		if v[1:] == "id" {
			result += strings.ToUpper(v[:1]) + "ID"
		} else {
			result += strings.ToUpper(v[:1]) + v[1:]
		}
	}
	return result
}

func getComment(decl *ast.GenDecl, typeSpec *ast.TypeSpec) string {
	comment := ""
	if decl.Doc != nil {
		comment = decl.Doc.Text()
		comment = strings.TrimSpace(strings.TrimPrefix(comment, typeSpec.Name.Name))
	}
	return comment
}

func getMethodComment(field *ast.Field) string {
	comment := ""
	if field.Doc != nil {
		comment = field.Doc.Text()
		comment = strings.TrimSpace(strings.TrimPrefix(comment, field.Names[0].Name))
	}
	return comment
}

func getMethodParams(field []*ast.Field) string {
	params := ""
	for _, param := range field {
		params += param.Names[0].Name + " " + getType(param.Type) + ", "
	}
	return strings.TrimSuffix(params, ", ")
}

func getMethodResult(field []*ast.Field) string {
	result := ""
	for _, param := range field {
		result += getType(param.Type) + ", "
	}
	if len(field) > 1 {
		return "(" + strings.TrimSuffix(result, ", ") + ")"
	}
	return strings.TrimSuffix(result, ", ")
}

func getMethodResultValues(field []*ast.Field) string {
	result := ""
	for _, param := range field {
		result += getTypeVal(param.Type) + ", "
	}
	return strings.TrimSuffix(result, ", ")
}

func getTypeVal(expr ast.Expr) string {
	switch expr.(type) {
	case *ast.Ident:
		switch expr.(*ast.Ident).Name {
		case "string":
			return `""`
		case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64",
			"float32", "float64", "complex64", "complex128":
			return "0"
		case "bool":
			return "false"
		}
		return "nil"
	case *ast.SelectorExpr:
		return "nil"
	case *ast.StarExpr:
		return "nil"
	case *ast.ArrayType:
		return "nil"
	case *ast.MapType:
		return "nil"
	default:
		return ""
	}
}

func getType(expr ast.Expr) string {
	switch expr.(type) {
	case *ast.Ident:
		return expr.(*ast.Ident).Name
	case *ast.SelectorExpr:
		return expr.(*ast.SelectorExpr).X.(*ast.Ident).Name + "." + expr.(*ast.SelectorExpr).Sel.Name
	case *ast.StarExpr:
		return "*" + getType(expr.(*ast.StarExpr).X)
	case *ast.ArrayType:
		return "[]" + getType(expr.(*ast.ArrayType).Elt)
	case *ast.MapType:
		return "map[" + getType(expr.(*ast.MapType).Key) + "]" + getType(expr.(*ast.MapType).Value)
	default:
		return ""
	}
}
