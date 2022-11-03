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

	"{apiPkg}"
)

type {ControllerName} struct {
}
`
const controllerFuncTpl = `
// {FuncName} {FuncDesc}
func ({simpleName} *{ControllerName}) {FuncName}(ctx context.Context, req *api.{ApiRequest}) (*api.{ApiResponse}, error) {

	return nil, errors.New("not implement")
}
`

// 生成controller
func (g *GenCmd) genController(apiFile string, decl *ast.GenDecl, specs *ast.TypeSpec) error {
	// 获取controller目录
	ctrlFile := filepath.Join(path.Dir(g.apiPath), "controller", strings.TrimPrefix(apiFile, g.apiPath))
	if err := os.MkdirAll(filepath.Dir(ctrlFile), 0755); err != nil {
		return err
	}
	// 生成controller
	_, err := os.Stat(ctrlFile)
	if err != nil {
		// 不存在重新生成
		if os.IsNotExist(err) {
			return g.regenController(ctrlFile, decl, apiFile, specs)
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
	funcTpl := g.genCtrlFunc(upperFirstChar(strings.TrimSuffix(filepath.Base(ctrlFile), ".go")), decl, specs)
	data = append(data, []byte(funcTpl)...)
	return os.WriteFile(ctrlFile, data, 0644)
}

// 重新生成controller
func (g *GenCmd) regenController(ctrlFile string, decl *ast.GenDecl, apiFile string, specs *ast.TypeSpec) error {
	// 生成controller头部
	data := controllerHeaderTpl
	ctrlName := upperFirstChar(strings.TrimSuffix(filepath.Base(ctrlFile), ".go"))
	data = strings.ReplaceAll(data, "{ControllerName}", ctrlName)
	abs, err := filepath.Abs(apiFile)
	if err != nil {
		return err
	}
	data = strings.ReplaceAll(data, "{apiPkg}", g.pkgName+strings.TrimPrefix(filepath.Dir(abs), g.pkgPath))

	data += g.genCtrlFunc(ctrlName, decl, specs)

	return os.WriteFile(ctrlFile, []byte(data), 0644)
}

func (g *GenCmd) genCtrlFunc(ctrlName string, decl *ast.GenDecl, specs *ast.TypeSpec) string {
	// 生成函数
	funcTpl := controllerFuncTpl
	funcTpl = strings.ReplaceAll(funcTpl, "{ControllerName}", ctrlName)
	funcTpl = strings.ReplaceAll(funcTpl, "{simpleName}", strings.ToLower(ctrlName[0:1]))
	funcName := strings.TrimSuffix(specs.Name.Name, "Request")
	funcTpl = strings.ReplaceAll(funcTpl, "{FuncName}", funcName)
	funcTpl = strings.ReplaceAll(funcTpl, "{ApiRequest}", specs.Name.Name)
	funcTpl = strings.ReplaceAll(funcTpl, "{ApiResponse}", funcName+"Response")
	desc := getComment(decl, specs)
	if desc == "" {
		desc = "在api中没有找到注释"
	}
	funcTpl = strings.ReplaceAll(funcTpl, "{FuncDesc}", desc)
	return funcTpl
}
