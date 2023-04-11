package definition

import (
	"fmt"
	"log"
	"testing"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"
)

func TestDefinition(t *testing.T) {

	src := `
		
	
	
	func test(){
		s2 := "s2"
		println s2
	}
	
	S := "test"
	println S
		
		

		
	`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	typeInfo := NewTypeInfo()
	typeInfo.analyze(file, fset)

	ast.Inspect(file, func(n ast.Node) bool {
		if ident, ok := n.(*ast.Ident); ok {
			def := typeInfo.findDefinition(ident)
			pos := fset.Position(ident.Pos())
			if def != nil {
				fmt.Printf("Found identifier '%s' at line %d, column %d with type '%v'\n", ident.Name, pos.Line, pos.Column, def)
			}
		}
		return true
	})
}
