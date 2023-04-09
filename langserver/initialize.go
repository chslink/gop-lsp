package langserver

import (
	"context"

	"gop-lsp/logger"
	lsp "gop-lsp/protocol"
	"gop-lsp/utils"
)

// InitializeParams 初始化参数
type InitializeParams struct {
	lsp.InitializeParams
	// InitializationOptions *InitializationOptions `json:"initializationOptions,omitempty"`
}

// Initialize lsp初始化函数
func (l *LspServer) Initialize(ctx context.Context, vs InitializeParams) (lsp.InitializeResult, error) {

	utils.InitialRootURIAndPath(string(vs.RootURI), vs.RootPath)
	logger.Debugf("Initialize ..., rootDir=%s, rootUri=%s", vs.RootPath, vs.RootURI)
	vscodeRoot := utils.VscodeURIToString(string(vs.RootURI))
	l.dir.SetVSRootDir(vscodeRoot)

	return lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{
			InnerServerCapabilities: lsp.InnerServerCapabilities{
				TextDocumentSync: &lsp.TextDocumentSyncOptions{
					OpenClose: true,
					Change:    lsp.Incremental,
					Save: lsp.SaveOptions{
						IncludeText: true,
					},
				},
				CompletionProvider: lsp.CompletionOptions{
					ResolveProvider: true,
					//TriggerCharacters: strings.Split("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.-:", ""),
					TriggerCharacters: []string{".", "\"", "'", ":", "-", "@", "#", " "},
					//AllCommitCharacters:[]string{".", "\"", "'", ":", "-", "@"},
				},
				ColorProvider: false,
				HoverProvider: true,
				// WorkspaceSymbolProvider: true,
				DefinitionProvider: true,
				ReferencesProvider: true,
				// DocumentSymbolProvider:  true,
				SignatureHelpProvider: lsp.SignatureHelpOptions{
					TriggerCharacters: []string{"(", ","},
				},
				CodeLensProvider: lsp.CodeLensOptions{
					ResolveProvider: false,
				},
				DocumentLinkProvider: lsp.DocumentLinkOptions{
					ResolveProvider: false,
				},
				RenameProvider:            true,
				DocumentHighlightProvider: true,
				Workspace: lsp.WorkspaceGn{
					WorkspaceFolders: lsp.WorkspaceFoldersGn{
						Supported:           true,
						ChangeNotifications: "workspace/didChangeWorkspaceFolders",
					},
				},
			},
		},
	}, nil
}

func (l *LspServer) initialCheckProject(ctx context.Context, workspaceFolderNum int, workspaceFolder []lsp.WorkspaceFolder) {
	l.dir.InitMainDir()

	for _, oneFloder := range workspaceFolder {
		logger.Debugf("floder=%s", oneFloder.URI)

		folderPath := utils.VscodeURIToString(oneFloder.URI)
		// 若增加的是当前workspace 文件夹中包含的子文件夹， 则不需要做任何处理
		if l.dir.IsDirExistWorkspace(folderPath) {
			logger.Debugf("current added dir=%s has existed in the workspaceFolder, not need analysis", folderPath)
			continue
		}
		l.dir.PushOneSubDir(folderPath)
	}

}

// InitializedParams 初始化参数
type InitializedParams struct {
	Settings interface{} `json:"settings"`
}

// Initialized 初始化
func (l *LspServer) Initialized(ctx context.Context, initialParam InitializedParams) error {
	logger.Debug("Initialized")
	// todo 获取所有的诊断错误
	// l.GetAllDiagnostics(ctx)
	return nil
}
