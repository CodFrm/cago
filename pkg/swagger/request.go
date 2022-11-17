package swagger

import (
	"go/ast"
	"go/parser"
	"go/token"
	"net/http"
	"strings"

	"github.com/codfrm/cago/internal/cmd/gen/utils"
	"github.com/go-openapi/spec"
)

// 解析文件生成swagger
func (s *Swagger) parseFile(filename string) error {
	f, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	for _, v := range f.Decls {
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
		// 解析结尾为Request的struct
		if strings.HasSuffix(typeSpec.Name.Name, "Response") {
			if err := newParseStruct(filename, s, f).parseStruct(typeSpec); err != nil {
				return err
			}
			continue
		}
		// 解析http.Route
		for _, field := range structSpec.Fields.List {
			expr, ok := field.Type.(*ast.SelectorExpr)
			if !ok {
				continue
			}
			if expr.Sel.Name != "Route" || expr.X.(*ast.Ident).Name != "http" {
				continue
			}
			// 处理swagger route
			if err := s.parseRoute(filename, f, decl, typeSpec, structSpec, field); err != nil {
				return err
			}
			// 处理swagger response
			break
		}
	}
	return nil
}

func (s *Swagger) parseRoute(filename string, file *ast.File, decl *ast.GenDecl,
	typeSpec *ast.TypeSpec, structSpec *ast.StructType, field *ast.Field) error {
	// 解析tag,取出route和method等参数
	tag := strings.TrimPrefix(field.Tag.Value, "`")
	tag = strings.TrimSuffix(tag, "`")
	// 取出path值
	path := utils.ParseTag(tag, "path")
	// 取出method值
	method := utils.ParseTag(tag, "method")
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
	text := utils.GetTypeComment(decl, typeSpec)
	operation.Summary = text
	operation.Description = text
	// 解析参数
	operation.Parameters = []spec.Parameter{}
	// GET请求参数放在query中
	contentType := utils.ParseTag(tag, "contentType")
	if contentType == "" {
		contentType = s.defaultContentType
	}
	operation.OperationProps.Consumes = []string{"application/" + contentType}
	// get请求,或者是非json请求,将参数单独解析出来
	if method == http.MethodGet || contentType != JSONBodyType {
		for _, field := range structSpec.Fields.List {
			// 解析参数
			if field.Names == nil {
				continue
			}
			name := utils.LowerFirstChar(field.Names[0].Name)
			tag := strings.TrimPrefix(field.Tag.Value, "`")
			in := "query"
			if method == http.MethodGet {
				in = "query"
			} else {
				uri := utils.ParseTag(tag, "uri")
				if uri != "" {
					in = "path"
				} else {
					in = "formData"
				}
			}
			validate := utils.ParseTag(tag, "validate")
			required := false
			if strings.Index(validate, "required") != -1 {
				required = true
			}
			if in == "path" {
				required = true
				path = strings.ReplaceAll(path, ":"+name, "{"+name+"}")
			}
			schema, err := newParseStruct(filename, s, file).parseFieldSwagger(field)
			if err != nil {
				return err
			}
			paramProps := spec.ParamProps{
				Description:     schema.Description,
				Name:            utils.LowerFirstChar(name),
				In:              in,
				Required:        required,
				AllowEmptyValue: false,
			}
			if schema.Type != nil {
				operation.Parameters = append(operation.Parameters, spec.Parameter{
					ParamProps: paramProps,
					SimpleSchema: spec.SimpleSchema{
						Type: schema.Type[0],
					},
				})
			} else {
				// 获取ref
				schema = s.swagger.Definitions[strings.Split(schema.Ref.GetURL().Fragment, "/")[2]]
				paramProps.Description += "\n" + schema.Description
				paramProps.Description = strings.TrimSpace(paramProps.Description)
				operation.Parameters = append(operation.Parameters, spec.Parameter{
					ParamProps: paramProps,
					CommonValidations: spec.CommonValidations{
						Enum: schema.Enum,
					},
					SimpleSchema: spec.SimpleSchema{
						Type: schema.Type[0],
					},
				})
			}
		}
	} else {
		// json请求,将参数放在body中
		schema, err := newParseStruct(filename, s, file).parseFieldType(typeSpec.Type)
		if err != nil {
			return err
		}
		operation.Parameters = append(operation.Parameters, spec.Parameter{
			ParamProps: spec.ParamProps{
				Description: schema.Description,
				Name:        "body",
				In:          "body",
			},
			SimpleSchema: spec.SimpleSchema{
				Type: schema.Type[0],
			},
		})
		//operation.Parameters = append(operation.Parameters, spec.Parameter{
		//	ParamProps: spec.ParamProps{
		//		Name: "body",
		//		In:   "body",
		//		Schema: &spec.Schema{
		//			SchemaProps: spec.SchemaProps{
		//				Ref: spec.MustCreateRef("#/definitions/" + file.Name.Name + "." + typeSpec.Name.Name),
		//			},
		//		},
		//	},
		//})
	}

	// 解析返回值
	operation.Responses.StatusCodeResponses = map[int]spec.Response{http.StatusOK: {
		ResponseProps: spec.ResponseProps{
			Description: "OK",
			Schema: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Ref: spec.MustCreateRef("#/definitions/" + file.Name.Name + "." +
						strings.Replace(typeSpec.Name.Name, "Request", "Response", 1)),
				},
			},
		},
	}, http.StatusBadRequest: {
		ResponseProps: spec.ResponseProps{
			Description: "Bad Request",
			Schema: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Ref: spec.MustCreateRef("#/definitions/BadRequest"),
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

	s.swagger.Paths.Paths[path] = pathItem
	return nil
}