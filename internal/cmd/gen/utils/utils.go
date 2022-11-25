package utils

import (
	"fmt"
	"go/ast"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// UpperFirstChar 首字母大写
func UpperFirstChar(str string) string {
	if len(str) == 0 {
		return str
	}
	return strings.ToUpper(str[0:1]) + str[1:]
}

// LowerFirstChar 首字母小写
func LowerFirstChar(str string) string {
	if len(str) == 0 {
		return str
	}
	return strings.ToLower(str[:1]) + str[1:]
}

// ToCamel 下划线转驼峰
func ToCamel(str string) string {
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

func GetTypeComment(decl *ast.GenDecl, typeSpec *ast.TypeSpec) string {
	comment := ""
	if decl.Doc != nil {
		comment = decl.Doc.Text()
		comment = strings.TrimSpace(strings.TrimPrefix(comment, typeSpec.Name.Name))
	}
	return comment
}

func GetFieldComment(field *ast.Field) string {
	comment := ""
	if field.Doc != nil {
		comment = field.Doc.Text()
		comment = strings.TrimSpace(strings.TrimPrefix(comment, field.Names[0].Name))
	} else if field.Comment != nil {
		comment = field.Comment.Text()
		comment = strings.TrimSpace(strings.TrimPrefix(comment, "//"))
	}
	return comment
}

func GetMethodParams(field []*ast.Field) string {
	params := ""
	for _, param := range field {
		params += param.Names[0].Name + " " + GetType(param.Type) + ", "
	}
	return strings.TrimSuffix(params, ", ")
}

func GetMethodResult(field []*ast.Field) string {
	result := ""
	for _, param := range field {
		result += GetType(param.Type) + ", "
	}
	if len(field) > 1 {
		return "(" + strings.TrimSuffix(result, ", ") + ")"
	}
	return strings.TrimSuffix(result, ", ")
}

func GetMethodResultValues(field []*ast.Field) string {
	result := ""
	for _, param := range field {
		result += GetTypeDefaultVal(param.Type) + ", "
	}
	return strings.TrimSuffix(result, ", ")
}

func GetTypeDefaultVal(expr ast.Expr) string {
	switch expr := expr.(type) {
	case *ast.Ident:
		switch expr.Name {
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

func GetType(expr ast.Expr) string {
	switch expr := expr.(type) {
	case *ast.Ident:
		return expr.Name
	case *ast.SelectorExpr:
		return expr.X.(*ast.Ident).Name + "." + expr.Sel.Name
	case *ast.StarExpr:
		return "*" + GetType(expr.X)
	case *ast.ArrayType:
		return "[]" + GetType(expr.Elt)
	case *ast.MapType:
		return "map[" + GetType(expr.Key) + "]" + GetType(expr.Value)
	default:
		return ""
	}
}

func ReadDir(path string, gen func(path string) error) error {
	dir, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, v := range dir {
		path := fmt.Sprintf("%s/%s", path, v.Name())
		if v.IsDir() {
			// 目录继续遍历
			if err := ReadDir(path, gen); err != nil {
				return err
			}
		} else {
			// 文件
			if !strings.HasSuffix(v.Name(), ".go") {
				continue
			}
			if err := gen(path); err != nil {
				return err
			}
		}
	}
	return nil
}

// FindRootPkgName 根据go.mod搜寻根包名
func FindRootPkgName(dir string) (string, string, error) {
	f, err := os.OpenFile(path.Join(dir, "./go.mod"), os.O_RDONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			// 向上层继续搜索
			absDir, err := filepath.Abs(dir)
			if err != nil {
				return "", "", err
			}
			return FindRootPkgName(path.Dir(absDir))
		}
		return "", "", err
	}
	defer f.Close()
	// 解析go.mod
	b, err := io.ReadAll(f)
	if err != nil {
		return "", "", err
	}
	// 解析出包名
	moduleName := strings.Split(string(b), "module ")[1]
	moduleName = strings.Split(moduleName, "\n")[0]
	pkgPath, err := filepath.Abs(dir)
	if err != nil {
		return "", "", err
	}
	return pkgPath, moduleName, nil
}

func PkgToPath(rootPkg, rootPkgName, pkgName string) (string, error) {
	// 去掉根包名
	pkgPath := strings.TrimPrefix(pkgName, rootPkgName)
	if pkgPath == pkgName {
		return "", fmt.Errorf("rootPkg must be prefix with rootPkgName")
	}
	// 拼接路径
	return path.Join(rootPkg, pkgPath), nil
}

func ParseTag(tag string, key string) string {
	keys := strings.Split(tag, key+":\"")
	if len(keys) != 2 {
		return ""
	}
	value := keys[1]
	value = strings.Split(value, "\"")[0]
	return value
}

func FileNameToCamel(filename string) string {
	return UpperFirstChar(ToCamel(strings.TrimSuffix(filepath.Base(filename), ".go")))
}
