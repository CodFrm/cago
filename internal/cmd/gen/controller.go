package gen

import (
	"go/ast"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const controllerHeaderTpl = `package controller

import (
	"context"
	"errors"

	"{ApiPkg}"
)

type {ControllerName} struct {
}

func New{ControllerName}() {ControllerName} {
	return {ControllerName}{}
}
`
const controllerFuncTpl = `
// {FuncName} {FuncDesc}
func ({SimpleName} *{ControllerName}) {FuncName}(ctx context.Context, req *{ApiPkgShort}.{ApiRequest}) (*{ApiPkgShort}.{ApiResponse}, error) {

	return nil, errors.New("not implement")
}
`

// 生成controller
func (c *Cmd) genController(apiFile string, f *ast.File, decl *ast.GenDecl, specs *ast.TypeSpec) error {
	// 获取controller目录
	ctrlFile := filepath.Join(path.Dir(c.apiPath), "controller", strings.TrimPrefix(apiFile, c.apiPath))
	if err := os.MkdirAll(filepath.Dir(ctrlFile), 0755); err != nil {
		return err
	}
	// 生成controller
	_, err := os.Stat(ctrlFile)
	if err != nil {
		// 不存在重新生成
		if os.IsNotExist(err) {
			return c.regenController(ctrlFile, f, decl, apiFile, specs)
		}
		return err
	}
	// 存在则判断是否需要添加新方法
	data, err := os.ReadFile(ctrlFile)
	if err != nil {
		return err
	}
	if strings.Contains(string(data), specs.Name.Name) {
		return nil
	}
	// 生成函数
	funcTpl := c.genCtrlFunc(upperFirstChar(strings.TrimSuffix(filepath.Base(ctrlFile), ".go")), f, decl, specs)
	data = append(data, []byte(funcTpl)...)
	return os.WriteFile(ctrlFile, data, 0644)
}

// 重新生成controller
func (c *Cmd) regenController(ctrlFile string, f *ast.File, decl *ast.GenDecl, apiFile string, specs *ast.TypeSpec) error {
	// 生成controller头部
	data := controllerHeaderTpl
	ctrlName := upperFirstChar(strings.TrimSuffix(filepath.Base(ctrlFile), ".go"))
	data = strings.ReplaceAll(data, "{ControllerName}", ctrlName)
	abs, err := filepath.Abs(apiFile)
	if err != nil {
		return err
	}
	data = strings.ReplaceAll(data, "{ApiPkg}", c.pkgName+strings.TrimPrefix(filepath.Dir(abs), c.pkgPath))

	data += c.genCtrlFunc(ctrlName, f, decl, specs)

	return os.WriteFile(ctrlFile, []byte(data), 0644)
}

func (c *Cmd) genCtrlFunc(ctrlName string, f *ast.File, decl *ast.GenDecl, specs *ast.TypeSpec) string {
	// 生成函数
	funcTpl := controllerFuncTpl
	funcTpl = strings.ReplaceAll(funcTpl, "{ControllerName}", ctrlName)
	funcTpl = strings.ReplaceAll(funcTpl, "{SimpleName}", strings.ToLower(ctrlName[0:1]))
	funcName := strings.TrimSuffix(specs.Name.Name, "Request")
	funcTpl = strings.ReplaceAll(funcTpl, "{FuncName}", funcName)
	funcTpl = strings.ReplaceAll(funcTpl, "{ApiRequest}", specs.Name.Name)
	funcTpl = strings.ReplaceAll(funcTpl, "{ApiResponse}", funcName+"Response")
	funcTpl = strings.ReplaceAll(funcTpl, "{ApiPkgShort}", f.Name.Name)
	desc := getComment(decl, specs)
	if desc == "" {
		desc = "在api中没有找到注释"
	}
	funcTpl = strings.ReplaceAll(funcTpl, "{FuncDesc}", desc)
	return funcTpl
}
