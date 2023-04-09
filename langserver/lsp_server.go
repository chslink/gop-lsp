package langserver

import (
	"context"
	"sync"
	"time"

	"gop-lsp/langserver/check"
	"gop-lsp/logger"
	lsp "gop-lsp/protocol"
	"gop-lsp/utils"

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

	dir *check.DirManager

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
		dir:            check.CreateDirManager(),
		colorTime:      0,
		changeConfFlag: false,
	}

	return lspServer
}

// CreateServer 创建server
func CreateServer() *jrpc2.Server {
	server := CreateLspServer()

	server.server = jrpc2.NewServer(handler.Map{
		"initialize":              handler.New(server.Initialize),
		"initialized":             handler.New(server.Initialized),
		"textDocument/didChange":  handler.New(server.TextDocumentDidChange),
		"textDocument/didSave":    handler.New(server.TextDocumentDidSave),
		"textDocument/didOpen":    handler.New(server.TextDocumentDidOpen),
		"textDocument/didClose":   handler.New(server.TextDocumentDidClose),
		"textDocument/definition": handler.New(server.TextDocumentDefine),
		// "textDocument/hover":             handler.New(server.TextDocumentHover),
		// "textDocument/references":        handler.New(server.TextDocumentReferences),
		// "textDocument/documentSymbol":    handler.New(server.TextDocumentSymbol),
		// "textDocument/rename":            handler.New(server.TextDocumentRename),
		// "textDocument/documentHighlight": handler.New(server.TextDocumentHighlight),
		// "textDocument/signatureHelp":     handler.New(server.TextDocumentSignatureHelp),
		// "textDocument/documentColor":     handler.New(server.TextDocumentColor),
		"textDocument/codeLens":     handler.New(server.TextDocumentCodeLens),
		"textDocument/documentLink": handler.New(server.TextDocumentLink),
		// "textDocument/completion":        handler.New(server.TextDocumentComplete),
		// "completionItem/resolve":         handler.New(server.TextDocumentCompleteResolve),
		// "workspace/didChangeConfiguration":    handler.New(server.ChangeConfiguration),
		// "workspace/didChangeWorkspaceFolders": handler.New(server.WorkspaceChangeWorkspaceFolders),
		// "workspace/didChangeWatchedFiles":     handler.New(server.WorkspaceChangeWatchedFiles),
		// "workspace/symbol":                    handler.New(server.WorkspaceSymbolRequest),
		// "luahelper/getVarColor":               handler.New(server.TextDocumentGetVarColor),
		// "luahelper/getOnlineReq":              handler.New(server.GetOnlineReq),
		"$/cancelRequest": handler.New(server.CancelRequest),
		"shutdown":        handler.New(server.Shutdown),
		"exit":            handler.New(server.Exit),
	}, &jrpc2.ServerOptions{
		AllowPush:   true,
		Concurrency: 4,
		Logger:      logger.Logger,
	})

	return server.server
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

// commFileRequest 通用的文件处理请求
type commFileRequest struct {
	result   bool         // 处理结果
	strFile  string       // 文件名
	contents []byte       // 文件具体的内容
	offset   int          // 光标的偏移
	pos      lsp.Position // 请求的行与列
}

// beginFileRequest 通用的文件处理请求预处理
func (l *LspServer) beginFileRequest(url lsp.DocumentURI, pos lsp.Position) (fileRequest commFileRequest) {
	fileRequest.result = false

	strFile := utils.VscodeURIToString(string(url))
	project := l.getAllProject()
	if !project.IsNeedHandle(strFile) {
		logger.Debug("not need to handle strFile=%s", strFile)
		return
	}

	fileCache := l.getFileCache()
	contents, found := fileCache.GetFileContent(strFile)
	if !found {
		logger.Printf("file %s not find contents", strFile)
		return
	}

	offset, err := utils.OffsetForPosition(contents, (int)(pos.Line), (int)(pos.Character))
	if err != nil {
		logger.Printf("file position error=%s", err.Error())
		return
	}

	fileRequest.result = true
	fileRequest.strFile = strFile
	fileRequest.contents = contents
	fileRequest.offset = offset
	fileRequest.pos = pos
	return
}
