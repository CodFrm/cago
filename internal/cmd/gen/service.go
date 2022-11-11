package gen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const serviceInterfaceTpl = `package {PkgName}

import (
	"context"

	api "{ApiPkg}"
)

type I{ServiceName} interface {
}
`

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

// 查找service文件生成service
func (c *Cmd) findService() error {
	serviceDir := path.Join(path.Dir(c.apiPath), "service")
	if err := os.MkdirAll(serviceDir, 0755); err != nil {
		return err
	}
	return c.readDir(serviceDir, func(path string) error {
		return c.genServiceMethod(path)
	})
}

func (c *Cmd) genService(apiFile string, f *ast.File, decl *ast.GenDecl, specs *ast.TypeSpec) error {
	// 生成service文件
	serviceFile := filepath.Join(path.Dir(c.apiPath), "service", strings.TrimPrefix(apiFile, c.apiPath))
	if err := os.MkdirAll(filepath.Dir(serviceFile), 0755); err != nil {
		return err
	}
	_, err := os.Stat(serviceFile)
	if err != nil {
		// 不存在重新生成
		if !os.IsNotExist(err) {
			return err
		}
		if err := c.regenService(serviceFile, f, apiFile); err != nil {
			return err
		}
	}
	data, err := os.ReadFile(serviceFile)
	if err != nil {
		return err
	}
	src := string(data)
	// 生成service接口方法
	serviceAst, err := parser.ParseFile(token.NewFileSet(), serviceFile, data, parser.ParseComments)
	if err != nil {
		return err
	}
	// 搜索接口
	for _, serviceDecl := range serviceAst.Decls {
		genDecl, ok := serviceDecl.(*ast.GenDecl)
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
		// 判断有没有service方法
		name := strings.TrimSuffix(specs.Name.Name, "Request")
		flag := false
		for _, method := range typeSpec.Type.(*ast.InterfaceType).Methods.List {
			if method.Names[0].Name == name {
				flag = true
				break
			}
		}
		if flag {
			return nil
		}
		// 插入service方法
		comment := getComment(decl, specs)
		if comment == "" {
			comment = "TODO"
		}
		data := "\t// " + name + " " + comment + "\n"
		data += "\t" +
			fmt.Sprintf("%s(ctx context.Context,req *api.%sRequest) (*api.%sResponse,error)\n",
				name,
				name, name)
		src = src[:serviceDecl.End()-2] + data + src[serviceDecl.End()-2:]
		return os.WriteFile(serviceFile, []byte(src), 0644)
	}
	return nil
}

func (c *Cmd) regenService(serviceFile string, f *ast.File, apiFile string) error {
	// 生成service头部
	data := serviceInterfaceTpl
	ctrlName := upperFirstChar(strings.TrimSuffix(filepath.Base(serviceFile), ".go"))
	data = strings.ReplaceAll(data, "{ServiceName}", ctrlName)
	data = strings.ReplaceAll(data, "{PkgName}", f.Name.Name)
	abs, err := filepath.Abs(apiFile)
	if err != nil {
		return err
	}
	data = strings.ReplaceAll(data, "{ApiPkg}", c.pkgName+strings.TrimPrefix(filepath.Dir(abs), c.pkgPath))
	return os.WriteFile(serviceFile, []byte(data), 0644)
}

func (c *Cmd) genServiceMethod(path string) error {
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
