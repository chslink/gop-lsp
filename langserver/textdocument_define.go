package langserver

import (
	"context"

	"gop-lsp/logger"
	lsp "gop-lsp/protocol"
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

	//strFile := fileRequest.strFile
	//project := l.getAllProject()

	// todo 1）判断查找的定义是否为打开一个文件 import package

	//

	return nil, err
}
