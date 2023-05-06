package definition

import (
	"fmt"
	"log"
	"testing"

	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"
)

func TestDefinition(t *testing.T) {

	src := `
	import "strings"

	func test(){
		s2 := "s2"
		s2 = strings.Split(s2,",")
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
	typeInfo.Analyze(file, fset)
	obj, err := typeInfo.FindDefinitionPos(116)
	fmt.Println(err)
	fmt.Println(obj)

}
