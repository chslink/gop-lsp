package langserver

import (
	"context"

	"gop-lsp/logger"
	lsp "gop-lsp/protocol"
)

// sendDiagnostics 给客户端推送错误诊断消息
func (l *LspServer) sendDiagnostics(ctx context.Context, diagnostics lsp.PublishDiagnosticsParams) {
	err := l.server.Notify(ctx, "textDocument/publishDiagnostics", diagnostics)
	if err != nil {
		logger.Debug("PushShowMessage error=%v", err)
	}
}
