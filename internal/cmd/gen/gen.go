package gen

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"strings"

	"github.com/codfrm/cago/internal/cmd/gen/utils"
	"github.com/codfrm/cago/pkg/swagger"
	"github.com/go-openapi/spec"
	"github.com/spf13/cobra"
)

const (
	JSONBodyType       = "json"
	FormDataBodyType   = "form-data"
	XWWWFormURLEncoded = "x-www-form-urlencoded"
)

type Cmd struct {
	apiPath     string
	pkgName     string
	pkgPath     string
	defaultBody string
}

func NewGenCmd() *Cmd {
	return &Cmd{}
}

func (c *Cmd) Commands() []*cobra.Command {
	ret := &cobra.Command{
		Use:   "gen",
		Short: "读取api目录下的文件,生成controller、service和swagger文档",
		RunE:  c.gen,
	}
	ret.AddCommand(&cobra.Command{
		Use:   "gorm [table]",
		Short: "输入表名,生成对应的model,需要配置好数据库连接",
		RunE:  c.genDB,
		Args:  cobra.ExactArgs(1),
	})
	ret.Flags().StringVarP(&c.apiPath, "dir", "d", "./internal/api", "api目录")
	return []*cobra.Command{ret}
}

func (c *Cmd) gen(cmd *cobra.Command, args []string) error {
	c.defaultBody = JSONBodyType
	var err error
	if err != nil {
		return err
	}
	c.pkgPath, c.pkgName, err = utils.FindRootPkgName(c.apiPath)
	if err != nil {
		return err
	}
	if err := utils.ReadDir(c.apiPath, func(path string) error {
		return c.genFile(path)
	}); err != nil {
		return err
	}
	// 生成swagger
	swagger := swagger.NewSwagger(c.apiPath)
	if err := swagger.Gen(); err != nil {
		return err
	}
	return swagger.Write()
}

func (c *Cmd) parseInfo() (*spec.Swagger, error) {
	// 解析main.go
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path.Join(c.pkgPath, "./main.go"), nil, parser.ParseComments)
	if err != nil {
		f, err = parser.ParseFile(fset, path.Join(c.pkgPath, "./cmd/app/main.go"), nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
	}
	info := &spec.Info{}
	ret := &spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Swagger: "2.0",
			Definitions: spec.Definitions{
				"BadRequest": {
					SchemaProps: spec.SchemaProps{
						Type: []string{"object"},
						Properties: map[string]spec.Schema{
							"code": {
								SchemaProps: spec.SchemaProps{
									Type:        []string{"integer"},
									Description: "错误码",
									Format:      "int32",
								},
							},
							"msg": {
								SchemaProps: spec.SchemaProps{
									Type:        []string{"string"},
									Description: "错误信息",
								},
							},
						},
					},
				},
			},
			Info: info,
			Paths: &spec.Paths{
				Paths: make(map[string]spec.PathItem),
			},
		},
	}
	for _, comment := range f.Comments {
		flag := false
		for _, v := range comment.List {
			text := strings.TrimPrefix(v.Text, "// @")
			// 证明是swagger的注释
			if text == v.Text {
				continue
			}
			flag = true
			// 解析注释
			key := strings.Split(text, " ")[0]
			value := strings.TrimPrefix(text, key+" ")
			value = strings.TrimSpace(value)
			switch strings.ToLower(key) {
			case "title":
				info.Title = value
			case "description":
				info.Description = value
			case "version":
				info.Version = value
			case "basepath":
				ret.BasePath = value
			}
		}
		if flag {
			break
		}
	}
	return ret, err
}

// 解析生成文件
func (c *Cmd) genFile(filepath string) error {
	// ast解析并生成swagger文档
	f, err := parser.ParseFile(token.NewFileSet(), filepath, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	for _, v := range f.Decls {
		// 解析带有Route的struct
		decl, ok := v.(*ast.GenDecl)
		if !ok {
			continue
		}
		if decl.Tok != token.TYPE {
			continue
		}
		typeSpec := decl.Specs[0].(*ast.TypeSpec)
		structSpec, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			continue
		}
		// 解析http.Route
		var routeField *ast.Field
		for _, field := range structSpec.Fields.List {
			expr, ok := field.Type.(*ast.SelectorExpr)
			if !ok {
				continue
			}
			if expr.Sel.Name != "Route" || expr.X.(*ast.Ident).Name != "mux" {
				continue
			}
			routeField = field
			break
		}
		if routeField == nil {
			continue
		}
		// 生成controller
		if err := c.genController(filepath, f, decl, typeSpec, routeField); err != nil {
			return err
		}
		// 生成service接口
		if err := c.genService(filepath, f, decl, typeSpec); err != nil {
			return err
		}
	}
	// 读取service目录根据接口生成service
	return c.findService()
}
