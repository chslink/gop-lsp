package langserver

import (
	"context"

	"github.com/yinfei8/jrpc2"
	"github.com/yinfei8/jrpc2/handler"
	"gop-lsp/logger"
	lsp "gop-lsp/protocol"
)

func createLspTest(strRootPath string, strRootUri string) *LspServer {
	//common.GlobalConfigDefautInit()
	//common.GConfig.IntialGlobalVar()

	server := CreateLspServer()
	server.server = jrpc2.NewServer(handler.Map{}, &jrpc2.ServerOptions{
		AllowPush:   false,
		Concurrency: 1,
	})

	ctx := context.Background()
	initializeParams := InitializeParams{
		InitializeParams: lsp.InitializeParams{
			InnerInitializeParams: lsp.InnerInitializeParams{
				RootPath: strRootPath,
				RootURI:  lsp.DocumentURI(strRootUri),
			},
		},
	}
	_, err := server.Initialize(ctx, initializeParams)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return server
}
