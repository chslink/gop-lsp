package langserver

import (
	"context"
	"gop-lsp/langserver/check"
	"gop-lsp/logger"
	"gop-lsp/utils"
	"sync"
	"time"

	"github.com/yinfei8/jrpc2"
	"github.com/yinfei8/jrpc2/handler"
)

type serverState int

const (
	serverCreated      = serverState(iota)
	serverInitializing //  initialize request
	serverInitialized  // initialized request
	serverShutDown
)

// LspServer lsp调用的全局对象
type LspServer struct {
	// 与客户端json rpc2通信的对象
	server *jrpc2.Server

	// 管理所有Gop工程的对象
	project *check.AllProject

	// 打开文件的缓冲
	fileCache *utils.FileMapCache

	// // 所有文件的诊断错误信息, 静态的
	// fileErrorMap map[string][]common.CheckError

	// // 所有文件的诊断错误信息, 动态的，文件实时修改了，但是没有保存的错误
	// fileChangeErrorMap map[string][]common.CheckError

	// 请求互斥锁
	requestMutex sync.Mutex

	// 最后一次获取文档着色功能的时间
	colorTime int64

	// 是否处理过ChangeConfiguration 标记
	changeConfFlag bool

	stateMu sync.Mutex
	state   serverState
}

// 管理指针
var lspServer *LspServer = nil

// CreateLspServer 创建全局Glsp管理对象
func CreateLspServer() *LspServer {
	lspServer = &LspServer{
		server:         nil,
		project:        nil,
		fileCache:      utils.CreateFileMapCache(),
		colorTime:      0,
		changeConfFlag: false,
	}

	return lspServer
}

// CreateServer 创建server
func CreateServer() *jrpc2.Server {
	lspServer := CreateLspServer()

	lspServer.server = jrpc2.NewServer(handler.Map{
		"initialize":             handler.New(lspServer.Initialize),
		"initialized":            handler.New(lspServer.Initialized),
		"textDocument/didChange": handler.New(lspServer.TextDocumentDidChange),
		"textDocument/didSave":   handler.New(lspServer.TextDocumentDidSave),
		"textDocument/didOpen":   handler.New(lspServer.TextDocumentDidOpen),
		"textDocument/didClose":  handler.New(lspServer.TextDocumentDidClose),
		// "textDocument/definition":        handler.New(lspServer.TextDocumentDefine),
		// "textDocument/hover":             handler.New(lspServer.TextDocumentHover),
		// "textDocument/references":        handler.New(lspServer.TextDocumentReferences),
		// "textDocument/documentSymbol":    handler.New(lspServer.TextDocumentSymbol),
		// "textDocument/rename":            handler.New(lspServer.TextDocumentRename),
		// "textDocument/documentHighlight": handler.New(lspServer.TextDocumentHighlight),
		// "textDocument/signatureHelp":     handler.New(lspServer.TextDocumentSignatureHelp),
		// "textDocument/documentColor":     handler.New(lspServer.TextDocumentColor),
		// "textDocument/codeLens":          handler.New(lspServer.TextDocumentCodeLens),
		// "textDocument/documentLink":      handler.New(lspServer.TextDocumentdocumentLink),
		// "textDocument/completion":        handler.New(lspServer.TextDocumentComplete),
		// "completionItem/resolve":         handler.New(lspServer.TextDocumentCompleteResolve),
		// "workspace/didChangeConfiguration":    handler.New(lspServer.ChangeConfiguration),
		// "workspace/didChangeWorkspaceFolders": handler.New(lspServer.WorkspaceChangeWorkspaceFolders),
		// "workspace/didChangeWatchedFiles":     handler.New(lspServer.WorkspaceChangeWatchedFiles),
		// "workspace/symbol":                    handler.New(lspServer.WorkspaceSymbolRequest),
		// "luahelper/getVarColor":               handler.New(lspServer.TextDocumentGetVarColor),
		// "luahelper/getOnlineReq":              handler.New(lspServer.GetOnlineReq),
		// "$/cancelRequest": handler.New(lspServer.CancelRequest),
		"shutdown": handler.New(lspServer.Shutdown),
		"exit":     handler.New(lspServer.Exit),
	}, &jrpc2.ServerOptions{
		AllowPush:   true,
		Concurrency: 4,
		Logger:      logger.Logger,
	})

	return lspServer.server
}

// getFileCache 获取文件缓冲map
func (l *LspServer) getFileCache() *utils.FileMapCache {
	return l.fileCache
}

// setColorTime 设置获取color着色的时间
func (l *LspServer) setColorTime(timeValue int64) {
	l.colorTime = timeValue
}

// isCanHighlight 判断是否可以对变量着色功能, 防止修改文件过程中频繁调用着色功能
func (l *LspServer) isCanHighlight() bool {
	// 如果修改文件的时间太频繁，返回false
	nowTime := time.Now().Unix()
	return nowTime-l.colorTime >= 3
}

// Shutdown lsp 关闭
func (l *LspServer) Shutdown(ctx context.Context) error {
	logger.Debug("Shutdown")
	return nil
}

// Exit 退出了
func (l *LspServer) Exit(ctx context.Context) error {
	logger.Debug("Exit")
	return nil
}

// getAllProject 获取CheckProject
func (l *LspServer) getAllProject() *check.AllProject {
	return l.project
}
