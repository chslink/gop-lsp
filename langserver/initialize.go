package langserver

import (
	"context"
	"gop-lsp/logger"
	lsp "gop-lsp/protocol"
)

// InitializeParams 初始化参数
type InitializeParams struct {
	lsp.InitializeParams
	// InitializationOptions *InitializationOptions `json:"initializationOptions,omitempty"`
}

// Initialize lsp初始化函数
func (l *LspServer) Initialize(ctx context.Context, vs InitializeParams) (lsp.InitializeResult, error) {
	logger.Debug("test init")
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

// InitializedParams 初始化参数
type InitializedParams struct {
	Settings interface{} `json:"settings"`
}

// Initialized 初始化
func (l *LspServer) Initialized(ctx context.Context, initialParam InitializedParams) error {
	logger.Debug("Initialized")
	// 获取所有的诊断错误
	// l.GetAllDiagnostics(ctx)
	return nil
}
