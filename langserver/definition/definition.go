package definition

import (
	"bytes"
	"go/printer"
	"path"
	"strings"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/token"
)

type GopObject struct {
	Name    string
	Type    string
	Pos     int
	Params  []string
	Results []string

	StructFields     []*GopObject
	InterfaceMethods []*GopObject
}

type ImportedPkg struct {
	Name string
	Path string
}

type TypeInfo struct {
	Defs         map[*ast.Ident]*GopObject
	ImportedPkgs map[string]*ImportedPkg
}

func NewTypeInfo() *TypeInfo {
	return &TypeInfo{
		Defs:         make(map[*ast.Ident]*GopObject),
		ImportedPkgs: map[string]*ImportedPkg{},
	}
}

func (ti *TypeInfo) findDefinition(ident *ast.Ident) *GopObject {
	if ti == nil || ident == nil {
		return nil
	}

	return ti.Defs[ident]
}

func (ti *TypeInfo) analyze(file *ast.File, fset *token.FileSet) {
	ast.Inspect(file, func(n ast.Node) bool {
		switch v := n.(type) {
		// Handle imports
		case *ast.ImportSpec:
			var pkgName string
			if v.Name != nil {
				pkgName = v.Name.Name
			} else {
				pkgName = path.Base(v.Path.Value)
				pkgName = strings.Trim(pkgName, `"`)
			}
			ti.ImportedPkgs[pkgName] = &ImportedPkg{
				Name: pkgName,
				Path: v.Path.Value,
			}

		// Handle variable assignments
		case *ast.AssignStmt:
			for i, expr := range v.Rhs {
				if ident, ok := v.Lhs[i].(*ast.Ident); ok {
					switch e := expr.(type) {
					case *ast.BasicLit:
						ti.Defs[ident] = &GopObject{Name: ident.Name, Type: e.Kind.String(), Pos: fset.Position(ident.Pos()).Offset}
					case *ast.Ident:
						if def, ok := ti.Defs[e]; ok {
							ti.Defs[ident] = &GopObject{Name: ident.Name, Type: def.Type, Pos: fset.Position(ident.Pos()).Offset}
						}
					}
				}
			}

		// Handle function declarations
		case *ast.FuncDecl:
			funcName := v.Name.Name
			// ... Add function handling logic here
			// Handle function parameters
			var params []string
			if v.Type.Params != nil {
				for _, field := range v.Type.Params.List {
					paramType := exprToString(field.Type)
					for _ = range field.Names {
						params = append(params, paramType)
					}
				}
			}

			// Handle function results (return values)
			var results []string
			if v.Type.Results != nil {
				for _, field := range v.Type.Results.List {
					resultType := exprToString(field.Type)
					for _ = range field.Names {
						results = append(results, resultType)
					}
				}
			}

			// Save function information in TypeInfo
			ti.Defs[v.Name] = &GopObject{
				Name:    funcName,
				Type:    "func",
				Pos:     fset.Position(v.Pos()).Offset,
				Params:  params,
				Results: results,
			}

		// Handle struct declarations
		case *ast.TypeSpec:
			typeName := v.Name.Name
			switch t := v.Type.(type) {
			case *ast.StructType:
				var fields []*GopObject
				for _, field := range t.Fields.List {
					fieldType := exprToString(field.Type)
					for _, fieldName := range field.Names {
						fields = append(fields, &GopObject{Name: fieldName.Name, Type: fieldType})
					}
				}

				ti.Defs[v.Name] = &GopObject{
					Name:         typeName,
					Type:         "struct",
					Pos:          fset.Position(v.Pos()).Offset,
					StructFields: fields,
				}

			case *ast.InterfaceType:
				var methods []*GopObject
				for _, method := range t.Methods.List {
					methodType := method.Type.(*ast.FuncType)

					// Handle method parameters
					var params []string
					if methodType.Params != nil {
						for _, field := range methodType.Params.List {
							paramType := exprToString(field.Type)
							for _ = range field.Names {
								params = append(params, paramType)
							}
						}
					}

					// Handle method results (return values)
					var results []string
					if methodType.Results != nil {
						for _, field := range methodType.Results.List {
							resultType := exprToString(field.Type)
							for _ = range field.Names {
								results = append(results, resultType)
							}
						}
					}

					for _, methodName := range method.Names {
						methods = append(methods, &GopObject{Name: methodName.Name, Type: "func", Params: params, Results: results})
					}
				}

				ti.Defs[v.Name] = &GopObject{
					Name:             typeName,
					Type:             "interface",
					Pos:              fset.Position(v.Pos()).Offset,
					InterfaceMethods: methods,
				}
			}
		}
		return true
	})
}

// Helper function to convert ast.Expr to a string representation of the type
func exprToString(expr ast.Expr) string {
	var buf bytes.Buffer
	printer.Fprint(&buf, token.NewFileSet(), expr)
	return buf.String()
}
