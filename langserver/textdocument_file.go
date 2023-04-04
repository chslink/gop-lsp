package langserver

import (
	"context"
	"gop-lsp/logger"
	lsp "gop-lsp/protocol"
	"gop-lsp/utils"
	"time"
)

// TextDocumentDidOpen 打开了一个文件的请求
func (l *LspServer) TextDocumentDidOpen(ctx context.Context, vs lsp.DidOpenTextDocumentParams) error {
	l.requestMutex.Lock()
	defer l.requestMutex.Unlock()

	// 判断打开的文件，是否是需要分析的文件
	strFile := utils.VscodeURIToString(string(vs.TextDocument.URI))
	// todo 判断忽略文件等等
	l.setColorTime(0)
	fileCache := l.getFileCache()
	fileCache.SetFileContent(strFile, utils.StringToBytes(vs.TextDocument.Text))
	// todo 代码文件分析
	return nil
}

// TextDocumentDidChange 单个文件的内容变化了
func (l *LspServer) TextDocumentDidChange(ctx context.Context, vs lsp.DidChangeTextDocumentParams) error {
	l.requestMutex.Lock()
	defer l.requestMutex.Unlock()

	// 判断打开的文件，是否是需要分析的文件
	strFile := utils.VscodeURIToString(string(vs.TextDocument.URI))
	fileCache := l.getFileCache()
	contents, found := fileCache.GetFileContent(strFile)
	if !found {
		logger.Printf("ApplyContentChanges get strFile=%s error", strFile)
		return nil
	}

	time1 := time.Now()
	changeContents, err := fileCache.ApplyContentChanges(strFile, contents, vs.ContentChanges)
	if err != nil {
		logger.Printf("ApplyContentChanges strFile=%s errInfo=%s", strFile, err.Error())
		return nil
	}
	ftime := time.Since(time1).Milliseconds()
	logger.Debug("TextDocumentDidChang ApplyContentChanges, strFile=%s, time=%d", strFile, ftime)

	fileCache.SetFileContent(strFile, changeContents)
	contents, found = fileCache.GetFileContent(strFile)
	if !found {
		logger.Printf("ApplyContentChanges get strFile=%s error", strFile)
		return nil
	}

	// todo 代码文件分析

	// 设置文件修改的时间
	l.setColorTime(time.Now().Unix())
	return nil
}

// TextDocumentDidClose 文件关闭了
func (l *LspServer) TextDocumentDidClose(ctx context.Context, vs lsp.DidCloseTextDocumentParams) error {
	l.requestMutex.Lock()
	defer l.requestMutex.Unlock()

	// 判断打开的文件，是否是需要分析的文件
	strFile := utils.VscodeURIToString(string(vs.TextDocument.URI))
	fileCache := l.getFileCache()
	fileCache.DelFileContent(strFile)
	return nil
}

// TextDocumentDidSave 文件的内容进行保存
func (l *LspServer) TextDocumentDidSave(ctx context.Context, vs lsp.DidSaveTextDocumentParams) error {
	l.requestMutex.Lock()
	defer l.requestMutex.Unlock()

	// 判断打开的文件，是否是需要分析的文件
	strFile := utils.VscodeURIToString(string(vs.TextDocument.URI))
	fileCache := l.getFileCache()
	fileCache.SetFileContent(strFile, utils.StringToBytes(*vs.Text))
	// todo 文件分析
	return nil
}
