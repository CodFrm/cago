package gen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/go-openapi/jsonreference"
	"github.com/go-openapi/spec"
	"github.com/spf13/cobra"
)

const (
	JSONBodyType       = "json"
	FormDataBodyType   = "form-data"
	XWWWFormURLEncoded = "x-www-form-urlencoded"
)

type GenCmd struct {
	apiPath     string
	pkgName     string
	pkgPath     string
	defaultBody string
	swagger     *spec.Swagger
}

func NewGenCmd() *GenCmd {
	return &GenCmd{}
}

func (g *GenCmd) Commands() []*cobra.Command {
	ret := &cobra.Command{
		Use:   "gen",
		Short: "读取目录下的文件,生成controller和swagger文档",
		RunE:  g.gen,
	}
	ret.Flags().StringVarP(&g.apiPath, "dir", "d", "./internal/api", "api目录")
	return []*cobra.Command{ret}
}

func (g *GenCmd) gen(cmd *cobra.Command, args []string) error {
	g.defaultBody = XWWWFormURLEncoded
	var err error
	g.swagger, err = g.parseInfo()
	if err != nil {
		return err
	}
	g.pkgName, err = g.findRootPkgName(g.apiPath)
	if err != nil {
		return err
	}
	if err := g.readDir(g.apiPath); err != nil {
		return err
	}
	// 生成swagger文档
	if err := os.MkdirAll("./docs", 0755); err != nil {
		return err
	}
	b, err := yaml.Marshal(g.swagger)
	if err != nil {
		return err
	}
	return os.WriteFile("./docs/swagger.yaml", b, 0644)
}

func (g *GenCmd) parseInfo() (*spec.Swagger, error) {
	// 解析main.go
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path.Join(g.pkgPath, "./main.go"), nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	info := &spec.Info{}
	ret := &spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Swagger:     "2.0",
			Definitions: make(spec.Definitions),
			Info:        info,
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

// 根据go.mod搜寻根包名
func (g *GenCmd) findRootPkgName(dir string) (string, error) {
	f, err := os.OpenFile(path.Join(dir, "./go.mod"), os.O_RDONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			// 向上层继续搜索
			absDir, err := filepath.Abs(dir)
			if err != nil {
				return "", err
			}
			return g.findRootPkgName(path.Dir(absDir))
		}
		return "", err
	}
	defer f.Close()
	// 解析go.mod
	b, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}
	// 解析出包名
	moduleName := strings.Split(string(b), "module ")[1]
	moduleName = strings.Split(moduleName, "\n")[0]
	g.pkgPath, err = filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	return moduleName, nil
}

func (g *GenCmd) readDir(path string) error {
	dir, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, v := range dir {
		path := fmt.Sprintf("%s/%s", path, v.Name())
		if v.IsDir() {
			// 目录继续遍历
			if err := g.readDir(path); err != nil {
				return err
			}
		} else {
			// 文件
			if !strings.HasSuffix(v.Name(), ".go") {
				continue
			}
			if err := g.genFile(path); err != nil {
				return err
			}
		}
	}
	return nil
}

// 解析生成文件
func (g *GenCmd) genFile(filepath string) error {
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
		structSpec := typeSpec.Type.(*ast.StructType)
		// 解析http.Route

		flag := false
		for _, field := range structSpec.Fields.List {
			expr, ok := field.Type.(*ast.SelectorExpr)
			if !ok {
				continue
			}
			if expr.Sel.Name != "Route" || expr.X.(*ast.Ident).Name != "http" {
				continue
			}
			flag = true
			// 处理swagger path
			if err := g.dealSwaggerPath(f, decl, typeSpec, structSpec, field); err != nil {
				return err
			}
			break
		}
		if !flag {
			// 处理swagger definitions
			if err := g.dealSwaggerDefinitions(f, typeSpec, structSpec); err != nil {
				return err
			}
			continue
		}
		// 生成controller
		if err := g.genController(filepath, decl, typeSpec); err != nil {
			return err
		}
	}
	return nil
}

// 处理swagger path
func (g *GenCmd) dealSwaggerPath(f *ast.File, decl *ast.GenDecl, typeSpec *ast.TypeSpec, structSpec *ast.StructType, field *ast.Field) error {
	// 解析tag,取出route和method等参数
	tag := strings.TrimPrefix(field.Tag.Value, "`")
	tag = strings.TrimSuffix(tag, "`")
	// 取出path值
	path := parseTag(tag, "path")
	// 取出method值
	method := parseTag(tag, "method")
	pathItem := spec.PathItem{
		PathItemProps: spec.PathItemProps{},
	}
	operation := &spec.Operation{
		OperationProps: spec.OperationProps{
			Produces:  []string{"application/json"},
			Responses: &spec.Responses{},
		},
	}
	// 解析注释
	text := getComment(decl, typeSpec)
	operation.Summary = text
	operation.Description = text
	// 解析参数
	operation.Parameters = []spec.Parameter{}
	// GET请求参数放在query中
	bodyType := parseTag(tag, "body")
	if bodyType == "" {
		bodyType = g.defaultBody
	}
	operation.OperationProps.Consumes = []string{"application/" + bodyType}
	if method == http.MethodGet || bodyType == XWWWFormURLEncoded {
		for _, field := range structSpec.Fields.List {
			// 解析参数
			if field.Names == nil {
				continue
			}
			name := lowerFirstChar(field.Names[0].Name)
			tag := strings.TrimPrefix(field.Tag.Value, "`")
			in := parseTag(tag, "in")
			if in == "" {
				if method == http.MethodGet {
					in = "query"
				} else {
					in = "formData"
				}
			}
			validate := parseTag(tag, "validate")
			required := false
			if strings.Index(validate, "required") != -1 {
				required = true
			}
			if in == "path" {
				required = true
				path = strings.ReplaceAll(path, ":"+name, "{"+name+"}")
			}
			schema, err := g.fieldType(f, field)
			if err != nil {
				return err
			}
			paramProps := spec.ParamProps{
				Description:     schema.Description,
				Name:            lowerFirstChar(name),
				In:              in,
				Required:        required,
				AllowEmptyValue: false,
			}
			operation.Parameters = append(operation.Parameters, spec.Parameter{
				ParamProps: paramProps,
				SimpleSchema: spec.SimpleSchema{
					Type: schema.Type[0],
				},
			})
		}
	}

	// 解析返回值
	operation.Responses.StatusCodeResponses = map[int]spec.Response{http.StatusOK: {
		ResponseProps: spec.ResponseProps{
			Description: "OK",
			Schema: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Ref: spec.MustCreateRef("#/definitions/" + f.Name.Name + "." +
						strings.Replace(typeSpec.Name.Name, "Request", "Response", 1)),
				},
			},
		},
	}}

	switch method {
	case http.MethodGet:
		pathItem.PathItemProps.Get = operation
	case http.MethodPost:
		pathItem.PathItemProps.Post = operation
	case http.MethodPut:
		pathItem.PathItemProps.Put = operation
	case http.MethodDelete:
		pathItem.PathItemProps.Delete = operation
	}

	g.swagger.Paths.Paths[path] = pathItem
	return nil
}

// 解析swagger definitions
func (g *GenCmd) dealSwaggerDefinitions(file *ast.File, specs *ast.TypeSpec, structSpec *ast.StructType) error {
	swaggerSchema := spec.Schema{
		SchemaProps: spec.SchemaProps{
			Properties: make(spec.SchemaProperties),
			Type:       spec.StringOrArray{"object"},
		},
	}
	// 处理参数
	var err error
	for _, field := range structSpec.Fields.List {
		if field.Names == nil {
			continue
		}
		// 为了防止死循环,先创建站位
		swaggerSchema.SchemaProps.Properties[lowerFirstChar(field.Names[0].Name)] = spec.Schema{}
		swaggerSchema.SchemaProps.Properties[lowerFirstChar(field.Names[0].Name)], err = g.fieldToSwagger(file, field)
		if err != nil {
			return err
		}
	}

	g.swagger.Definitions[fmt.Sprintf("%s.%s", file.Name.Name, specs.Name.Name)] = swaggerSchema

	return nil
}

func (g *GenCmd) fieldToSwagger(file *ast.File, field *ast.Field) (spec.Schema, error) {
	ret, err := g.fieldType(file, field)
	if err != nil {
		return spec.Schema{}, err
	}
	description := ""
	if field.Doc != nil {
		description = strings.TrimPrefix(field.Doc.List[0].Text, "//")
		description = strings.TrimSpace(description)
		description = strings.TrimPrefix(description, field.Names[0].Name)
		description = strings.TrimSpace(description)
	} else if field.Comment != nil {
		description = strings.TrimPrefix(field.Comment.List[0].Text, "//")
		description = strings.TrimSpace(description)
	}
	ret.SchemaProps.Description = description
	return ret, nil
}

func (g *GenCmd) fieldType(file *ast.File, field *ast.Field) (spec.Schema, error) {
	var swaggerType spec.SchemaProps
	// 转换类型
	var typeName string
	var fieldType ast.Expr
	var isSelectorExpr bool
	arrayType, isArray := field.Type.(*ast.ArrayType)
	if isArray {
		fieldType = arrayType.Elt
		typeName = "array"
	} else {
		fieldType = field.Type
	}
	var selectorExpr *ast.SelectorExpr
	selectorExpr, isSelectorExpr = fieldType.(*ast.SelectorExpr)
	if isSelectorExpr {
		fieldType = selectorExpr.Sel
		typeName = "object"
	} else {
		typeName = fieldType.(*ast.Ident).Name
	}

	switch typeName {
	case "string":
		swaggerType.Type = spec.StringOrArray{"string"}
	case "int", "int64", "int32", "int16", "int8", "uint", "uint64", "uint32", "uint16", "uint8":
		swaggerType.Type = spec.StringOrArray{"integer"}
	case "float32", "float64":
		swaggerType.Type = spec.StringOrArray{"number"}
	case "bool":
		swaggerType.Type = spec.StringOrArray{"boolean"}
	case "object":
		// 解析嵌套类型
		var err error
		swaggerType, err = g.parseStruct(file, selectorExpr)
		if err != nil {
			return spec.Schema{}, err
		}
	}
	ret := spec.Schema{}
	if isArray {
		ret.SchemaProps.Type = spec.StringOrArray{"array"}
		ret.SchemaProps.Items = &spec.SchemaOrArray{
			Schema: &spec.Schema{
				SchemaProps: swaggerType,
			},
		}
	} else {
		ret.SchemaProps = swaggerType
	}

	return ret, nil
}

// 解析嵌套结构体
func (g *GenCmd) parseStruct(f *ast.File, selectorExpr *ast.SelectorExpr) (spec.SchemaProps, error) {
	// 先检查是否已有结构体
	pkgName := selectorExpr.X.(*ast.Ident).Name
	name := fmt.Sprintf("%s.%s", pkgName, selectorExpr.Sel.Name)
	ref := jsonreference.MustCreateRef("#/definitions/" + name)
	if _, ok := g.swagger.Definitions[name]; ok {
		return spec.SchemaProps{
			Ref: spec.Ref{Ref: ref},
		}, nil
	}
	// 找到文件并解析
	for _, v := range f.Imports {
		// 获取名称
		name := strings.TrimSuffix(v.Path.Value, `"`)
		name = path.Base(name)
		fullName := name
		if v.Name != nil {
			name = v.Name.Name
		}
		if name != pkgName {
			continue
		}
		// 将包名换为完整包名
		if name != fullName {
			pkgName = fullName
			ref = jsonreference.MustCreateRef("#/definitions/" + fmt.Sprintf("%s.%s", pkgName, selectorExpr.Sel.Name))
		}
		path := strings.TrimSuffix(v.Path.Value, `"`)
		path = strings.TrimPrefix(path, `"`)
		// 暂时不处理非本包的结构
		if !strings.HasPrefix(path, g.pkgName) {
			continue
		}
		// 转换成绝对路径
		path = filepath.Join(g.pkgPath, strings.TrimPrefix(path, g.pkgName))
		// 解析文件
		pkgs, err := parser.ParseDir(token.NewFileSet(), path, nil, parser.ParseComments)
		if err != nil {
			return spec.SchemaProps{}, err
		}
		for _, pkg := range pkgs {
			for _, file := range pkg.Files {
				obj, ok := file.Scope.Objects[selectorExpr.Sel.Name]
				if !ok {
					continue
				}
				typeSpec, ok := obj.Decl.(*ast.TypeSpec)
				if !ok {
					continue
				}
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}
				if err := g.dealSwaggerDefinitions(file, typeSpec, structType); err != nil {
					return spec.SchemaProps{}, err
				}
			}
		}
	}

	return spec.SchemaProps{
		Ref: spec.Ref{Ref: ref},
	}, nil
}

func parseTag(tag string, key string) string {
	keys := strings.Split(tag, key+":\"")
	if len(keys) != 2 {
		return ""
	}
	value := keys[1]
	value = strings.Split(value, "\"")[0]
	return value
}
