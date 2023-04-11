package langserver

import (
	"context"
	"fmt"

	"gop-lsp/logger"
	lsp "gop-lsp/protocol"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"
)

func (l *LspServer) TextDocumentDefine(ctx context.Context, vs lsp.TextDocumentPositionParams) (locList []lsp.Location, err error) {
	l.requestMutex.Lock()
	defer l.requestMutex.Unlock()

	fileRequest := l.beginFileRequest(vs.TextDocument.URI, vs.Position)
	if !fileRequest.result {
		logger.Printf("TextDocumentDefine beginFileRequest false, uri=%s", vs.TextDocument.URI)
		return
	}
	if len(fileRequest.contents) == 0 || fileRequest.offset >= len(fileRequest.contents) {
		return
	}

	strFile := fileRequest.strFile
	//project := l.getAllProject()

	// TODO: 解析源文件
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, strFile, fileRequest.contents, parser.AllErrors)
	if err != nil {
		return
	}
	// 3. 在 AST 中查找符号
	pos := fset.Position(token.Pos(fileRequest.offset))
	ident, _ := findIdentAtPosition(file, pos)

	// 4. 查找符号定义 (具体实现可能需要调整)
	obj := findDefinition(ident)
	if obj == nil {
		return nil, fmt.Errorf("symbol not found")
	}

	// 5. 构造响应
	definitionPos := fset.Position(obj.Pos())
	location := lsp.Location{
		URI: vs.TextDocument.URI,
		Range: lsp.Range{
			Start: lsp.Position{Line: uint32(definitionPos.Line - 1), Character: uint32(definitionPos.Column)},
			End:   lsp.Position{Line: uint32(definitionPos.Line - 1), Character: uint32(definitionPos.Column + len(obj.Name))},
		},
	}

	return []lsp.Location{location}, err
}

func findIdentAtPosition(file *ast.File, pos token.Position) (*ast.Ident, bool) {

	var foundIdent *ast.Ident

	ast.Inspect(file, func(n ast.Node) bool {
		if n == nil {
			return false
		}

		if ident, ok := n.(*ast.Ident); ok {
			identPos := ident.Pos()
			identEnd := ident.End()
			if posIsValid(pos, identPos, identEnd) {
				foundIdent = ident
				return false
			}
		}
		return true

	})
	if foundIdent != nil {
		return foundIdent, true
	}
	return nil, false
}

func posIsValid(pos token.Position, startPos, endPos token.Pos) bool {
	if !startPos.IsValid() || !endPos.IsValid() {
		return false
	}
	posOffset := pos.Offset
	startOffset := int(startPos)
	endOffset := int(endPos)
	return posOffset >= startOffset && posOffset < endOffset
}

func findDefinition(ident *ast.Ident) *ast.Object {
	// TODO: 在这里实现查找符号定义的逻辑
	// ...

	return nil
}
