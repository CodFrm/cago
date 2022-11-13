package swagger

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"strings"

	"github.com/codfrm/cago/internal/cmd/gen/utils"
	"github.com/go-openapi/jsonreference"
	"github.com/go-openapi/spec"
)

type parseStruct struct {
	filename string
	*Swagger
	f *ast.File
}

func newParseStruct(filename string, s *Swagger, f *ast.File) *parseStruct {
	return &parseStruct{
		filename: filename,
		Swagger:  s,
		f:        f,
	}
}

func (p *parseStruct) parseStruct(typeSpec *ast.TypeSpec) error {
	name := fmt.Sprintf("%s.%s", p.f.Name, typeSpec.Name.Name)
	// 判断是否生成过,生成过则跳过
	if _, ok := p.swagger.Definitions[name]; ok {
		return nil
	}
	schema, err := p.parseFieldType(typeSpec.Type)
	if err != nil {
		return err
	}
	p.swagger.Definitions[name] = schema
	return nil
}

func (p *parseStruct) parseFieldSwagger(field *ast.Field) (spec.Schema, error) {
	// 数组类型
	if expr, ok := field.Type.(*ast.ArrayType); ok {
		// 解析数组类型
		schema, err := p.parseFieldType(expr.Elt)
		if err != nil {
			return spec.Schema{}, err
		}
		return spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"array"},
				Items: &spec.SchemaOrArray{
					Schema: &schema,
				},
			},
		}, nil
	}
	schema, err := p.parseFieldType(field.Type)
	if err != nil {
		return spec.Schema{}, err
	}
	schema.Description = utils.GetFieldComment(field)
	return schema, nil
}

func (p *parseStruct) parseFieldType(fieldType ast.Expr) (spec.Schema, error) {
	var swaggerType spec.SchemaProps
	t, ok := fieldType.(*ast.Ident)
	if !ok {
		// 判断interface
		if _, ok := fieldType.(*ast.InterfaceType); ok {
			return spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: []string{"object"},
				},
			}, nil
		} else if mapType, ok := fieldType.(*ast.MapType); ok {
			schema, err := p.parseFieldType(mapType.Value)
			if err != nil {
				return spec.Schema{}, err
			}
			return spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: []string{"object"},
					AdditionalProperties: &spec.SchemaOrBool{
						Allows: true,
						Schema: &schema,
					},
				},
			}, nil
		} else if structType, ok := fieldType.(*ast.StructType); ok {
			// 内联结构体
			swaggerType.Properties = make(map[string]spec.Schema)
			for _, field := range structType.Fields.List {
				if field.Names == nil {
					continue
				}
				// 解析字段
				schema, err := p.parseFieldSwagger(field)
				if err != nil {
					return spec.Schema{}, err
				}
				swaggerType.Properties[utils.LowerFirstChar(field.Names[0].Name)] = schema
			}
			swaggerType.Type = []string{"object"}
			return spec.Schema{
				SchemaProps: swaggerType,
			}, nil
		}
		return p.parseExpr(fieldType)
	}
	typeName := t.Name
	switch typeName {
	case "string":
		swaggerType.Type = spec.StringOrArray{"string"}
	case "int", "int64", "int32", "int16", "int8", "uint", "uint64", "uint32", "uint16", "uint8":
		swaggerType.Type = spec.StringOrArray{"integer"}
	case "float32", "float64":
		swaggerType.Type = spec.StringOrArray{"number"}
	case "bool":
		swaggerType.Type = spec.StringOrArray{"boolean"}
	default:
		return p.parseExpr(t)
	}
	return spec.Schema{
		SchemaProps: swaggerType,
	}, nil
}

// 解析引用类型
func (p *parseStruct) parseExpr(expr ast.Expr) (spec.Schema, error) {
	var pkgName, structName string
	if selectorExpr, ok := expr.(*ast.SelectorExpr); ok {
		pkgName = selectorExpr.X.(*ast.Ident).Name
		structName = selectorExpr.Sel.Name
	} else if startExpr, ok := expr.(*ast.StarExpr); ok {
		return p.parseExpr(startExpr.X)
	} else if ident, ok := expr.(*ast.Ident); ok {
		pkgName = p.f.Name.Name
		structName = ident.Name
	}
	ref := fmt.Sprintf("#/definitions/%s.%s", pkgName, structName)
	return spec.Schema{
		SchemaProps: spec.SchemaProps{
			Ref: spec.Ref{
				Ref: jsonreference.MustCreateRef(ref),
			},
		},
	}, p.findStruct(pkgName, structName)
}

// 查找包文件并解析
func (p *parseStruct) findStruct(pkgName string, structName string) error {
	// 查找包文件
	for _, f := range p.f.Imports {
		dir := strings.Trim(f.Path.Value, "\"")
		if f.Name != nil && f.Name.Name == pkgName {
			return p.parseFile(dir)
		}
		if path.Base(dir) == pkgName {
			// 解析包文件,转化为文件路径
			dir, err := utils.PkgToPath(p.rootPkgPath, p.rootPkgName, dir)
			if err != nil {
				return err
			}
			return p.parseDir(dir, structName)
		}
	}
	// 未找到,并且包名相等
	if pkgName == p.f.Name.Name {
		// 同目录
		return p.parseDir(path.Dir(p.filename), structName)
	}
	return errors.New("not found")
}

// 解析指定目录下的指定类型
func (p *parseStruct) parseDir(dir string, structName string) error {
	pkgs, err := parser.ParseDir(token.NewFileSet(), dir, func(info os.FileInfo) bool {
		return true
	}, parser.ParseComments)
	if err != nil {
		return err
	}
	// 指定结构体
	for _, pkg := range pkgs {
		for filename, f := range pkg.Files {
			for _, decl := range f.Decls {
				if genDecl, ok := decl.(*ast.GenDecl); ok {
					for _, spec := range genDecl.Specs {
						if typeSpec, ok := spec.(*ast.TypeSpec); ok {
							if typeSpec.Name.Name == structName {
								return newParseStruct(path.Join(dir, path.Base(filename)), p.Swagger, f).
									parseStruct(typeSpec)
							}
						}
					}
				}
			}
		}
	}
	return errors.New("not found")
}
