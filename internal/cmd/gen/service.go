package gen

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"strings"
)

const serviceStructTpl = `
type {ServiceName} struct {
}

var default{UpperServiceName} = &{ServiceName}{}

func {UpperServiceName}() {ServiceInterface} {
	return default{UpperServiceName}
}
`

const serviceMethodTpl = `
// {MethodName} {Comment}
func ({FirstServiceName} *{ServiceName}) {MethodName}({MethodParams}) {MethodResult} {
	return {MethodResultValues}
}
`

func (c *Cmd) findService() error {
	serviceDir := path.Join(path.Dir(c.apiPath), "service")
	if err := os.MkdirAll(serviceDir, 0755); err != nil {
		return err
	}
	return c.readDir(serviceDir, func(path string) error {
		return c.genService(path)
	})
}

func (c *Cmd) genService(path string) error {
	f, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	// 搜索接口
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		typeSpec, ok := genDecl.Specs[0].(*ast.TypeSpec)
		if !ok {
			continue
		}
		if !strings.HasPrefix(typeSpec.Name.Name, "I") {
			continue
		}
		// 生成service
		if err := c.genServiceFile(path, f, genDecl, typeSpec); err != nil {
			return err
		}
		break
	}
	return nil
}

func (c *Cmd) genServiceFile(path string, f *ast.File, genDecl *ast.GenDecl, typeSpec *ast.TypeSpec) error {
	// 判断是否已经生成struct
	hasStruct := false
	serviceName := lowerFirstChar(strings.TrimPrefix(typeSpec.Name.Name, "I"))
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		sTypeSpec, ok := genDecl.Specs[0].(*ast.TypeSpec)
		if !ok {
			continue
		}
		if sTypeSpec.Name.Name == serviceName {
			hasStruct = true
			break
		}
	}
	appendStr := ""
	// 生成struct
	if !hasStruct {
		appendStr = serviceStructTpl
		appendStr = strings.ReplaceAll(appendStr, "{ServiceName}", serviceName)
		appendStr = strings.ReplaceAll(appendStr, "{UpperServiceName}", upperFirstChar(serviceName))
		appendStr = strings.ReplaceAll(appendStr, "{ServiceInterface}", typeSpec.Name.Name)
	}
	// 生成方法
	for _, method := range typeSpec.Type.(*ast.InterfaceType).Methods.List {
		// 判断是否已经生成
		hasMethod := false
		for _, decl := range f.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			if funcDecl.Name.Name == method.Names[0].Name {
				hasMethod = true
				break
			}
		}
		if hasMethod {
			continue
		}
		methodStr := serviceMethodTpl
		methodStr = strings.ReplaceAll(methodStr, "{MethodName}", method.Names[0].Name)
		methodStr = strings.ReplaceAll(methodStr, "{Comment}", getMethodComment(method))
		methodStr = strings.ReplaceAll(methodStr, "{FirstServiceName}", lowerFirstChar(serviceName)[0:1])
		methodStr = strings.ReplaceAll(methodStr, "{ServiceName}", serviceName)
		methodStr = strings.ReplaceAll(methodStr, "{MethodParams}", getMethodParams(method.Type.(*ast.FuncType).Params.List))
		methodStr = strings.ReplaceAll(methodStr, "{MethodResult}", getMethodResult(method.Type.(*ast.FuncType).Results.List))
		methodStr = strings.ReplaceAll(methodStr, "{MethodResultValues}", getMethodResultValues(method.Type.(*ast.FuncType).Results.List))
		appendStr += methodStr
	}
	// 写入文件
	w, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = w.WriteString(appendStr)
	return err
}

func (c *Cmd) genServiceMethod() {

}
