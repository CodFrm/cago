package gen

import (
	"go/ast"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/codfrm/cago/internal/cmd/gen/utils"
)

const controllerHeaderTpl = `package {PkgName}

import (
	{ContextPkg}

	api "{ApiPkg}"
	service "{ServicePkg}"
)

type {ControllerName} struct {
}

func New{ControllerName}() *{ControllerName} {
	return &{ControllerName}{}
}
`

const controllerFuncTpl = `
// {FuncName} {FuncDesc}
func ({SimpleName} *{ControllerName}) {FuncName}(ctx {Context}, req *api.{ApiRequest}) (*api.{ApiResponse}, error) {
	return service.{ControllerName}().{FuncName}({ContextParam}, req)
}
`

type controller struct {
}

// 生成controller
func (c *Cmd) genController(apiFile string, f *ast.File, decl *ast.GenDecl, specs *ast.TypeSpec, routeField *ast.Field) error {
	// 获取controller目录
	ctrlFile := filepath.Join(path.Dir(c.apiPath), "controller", strings.TrimPrefix(apiFile, c.apiPath))
	if err := os.MkdirAll(filepath.Dir(ctrlFile), 0755); err != nil {
		return err
	}
	ginContext := false
	if utils.ParseTag(routeField.Tag.Value, "context") == "gin" {
		ginContext = true
	}
	// 生成controller
	_, err := os.Stat(ctrlFile)
	if err != nil {
		// 不存在重新生成
		if os.IsNotExist(err) {
			return c.regenController(ctrlFile, f, decl, apiFile, specs, ginContext)
		}
		return err
	}
	// 存在则判断是否需要添加新方法
	data, err := os.ReadFile(ctrlFile)
	if err != nil {
		return err
	}
	if strings.Contains(string(data), strings.TrimSuffix(specs.Name.Name, "Request")) {
		return nil
	}
	// 生成函数
	funcTpl := c.genCtrlFunc(ctrlFile,
		decl, specs, ginContext)
	data = append(data, []byte(funcTpl)...)
	return os.WriteFile(ctrlFile, data, 0644)
}

// 重新生成controller
func (c *Cmd) regenController(ctrlFile string, f *ast.File, decl *ast.GenDecl,
	apiFile string, specs *ast.TypeSpec, ginContext bool) error {
	// 生成controller头部
	data := controllerHeaderTpl
	ctrlName := utils.FileNameToCamel(ctrlFile)
	data = strings.ReplaceAll(data, "{ControllerName}", ctrlName)
	data = strings.ReplaceAll(data, "{PkgName}", f.Name.Name)
	abs, err := filepath.Abs(apiFile)
	if err != nil {
		return err
	}
	if ginContext {
		data = strings.ReplaceAll(data, "{ContextPkg}", `"github.com/gin-gonic/gin"`)
	} else {
		data = strings.ReplaceAll(data, "{ContextPkg}", `"context"`)
	}
	data = strings.ReplaceAll(data, "{ApiPkg}", c.pkgName+strings.TrimPrefix(filepath.Dir(abs), c.pkgPath))
	// 获取service包名
	abs, err = filepath.Abs(c.apiPath)
	if err != nil {
		return err
	}
	servicePkg := c.pkgName + strings.TrimPrefix(filepath.Dir(abs), c.pkgPath) + "/service/" + strings.TrimPrefix(filepath.Dir(apiFile), "internal/api/")
	data = strings.ReplaceAll(data, "{ServicePkg}", servicePkg)

	log.Printf("生成controller: %s", ctrlName)

	data += c.genCtrlFunc(ctrlFile, decl, specs, ginContext)

	return os.WriteFile(ctrlFile, []byte(data), 0644)
}

func (c *Cmd) genCtrlFunc(ctrlFile string, decl *ast.GenDecl, specs *ast.TypeSpec, ginContext bool) string {
	// 生成函数
	funcTpl := controllerFuncTpl
	ctrlName := utils.FileNameToCamel(ctrlFile)
	funcTpl = strings.ReplaceAll(funcTpl, "{ControllerName}", ctrlName)
	funcTpl = strings.ReplaceAll(funcTpl, "{SimpleName}", strings.ToLower(ctrlName[0:1]))
	funcName := strings.TrimSuffix(specs.Name.Name, "Request")
	funcTpl = strings.ReplaceAll(funcTpl, "{FuncName}", funcName)
	funcTpl = strings.ReplaceAll(funcTpl, "{ApiRequest}", specs.Name.Name)
	funcTpl = strings.ReplaceAll(funcTpl, "{ApiResponse}", funcName+"Response")
	desc := utils.GetTypeComment(decl, specs)
	if desc == "" {
		desc = "TODO"
	}
	if ginContext {
		funcTpl = strings.ReplaceAll(funcTpl, "{Context}", "*gin.Context")
		funcTpl = strings.ReplaceAll(funcTpl, "{ContextParam}", "ctx.Request.Context()")
	} else {
		funcTpl = strings.ReplaceAll(funcTpl, "{Context}", "context.Context")
		funcTpl = strings.ReplaceAll(funcTpl, "{ContextParam}", "ctx")
	}
	funcTpl = strings.ReplaceAll(funcTpl, "{FuncDesc}", desc)
	log.Printf("生成controller函数: %s", funcName)
	return funcTpl
}
