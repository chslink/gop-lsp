package definition

import (
	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"
	"github.com/goplus/gop/x/types"
)

type DefinitionFinder struct {
	fset  *token.FileSet
	pkgs  map[string]*types.Package
	types *TypeInfo
}

func NewDefinitionFinder() *DefinitionFinder {
	return &DefinitionFinder{
		fset:  token.NewFileSet(),
		pkgs:  make(map[string]*types.Package),
		types: NewTypeInfo(),
	}
}

func (df *DefinitionFinder) findDefinition(ident *ast.Ident) *GopObject {
	if df.types == nil || ident == nil {
		return nil
	}
	return df.types.Defs[ident]
}

func (df *DefinitionFinder) loadPackage(packagePath string) error {
	packageAst, err := parser.ParseDir(df.fset, packagePath, nil, parser.AllErrors)
	if err != nil {
		return err
	}

	typeChecker := types.Config{}
	for _, pkgAst := range packageAst {
		pkg, err := types.Load(df.fset, pkgAst, &typeChecker)
		if err != nil {
			return err
		}
		df.pkgs[pkg.Name()] = pkg
	}

	return nil
}
