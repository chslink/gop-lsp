package langserver

import (
	"context"

	"gop-lsp/logger"
	lsp "gop-lsp/protocol"
)

// CancelRequest 取消一个请求
func (l *LspServer) CancelRequest(ctx context.Context, vs lsp.CancelParams) error {
	logger.Debug("CancelRequest, id=%v", vs.ID)
	return nil
}

// TextDocumentCodeLens 请求
func (l *LspServer) TextDocumentCodeLens(ctx context.Context, vs lsp.CodeLensParams) (edit []lsp.CodeLens, err error) {
	return
}

// TextDocumentLink 请求
func (l *LspServer) TextDocumentLink(ctx context.Context, vs lsp.DocumentLinkParams) (edit []lsp.DocumentLink, err error) {
	return
}
