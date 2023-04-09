package main

import (
	"flag"
	"io"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sync"

	"gop-lsp/langserver"
	"gop-lsp/logger"

	"github.com/yinfei8/jrpc2/channel"
)

func main() {
	modeFlag := flag.Int("mode", 0, "mode type, 0 is run cmd, 1 is local rpc, 2 is socket rpc")
	logFlag := flag.Int("logflag", 0, "0 is not open log, 1 is open log")
	// localpath := flag.String("localpath", "", "local project path")
	flag.Parse()

	// 是否开启日志
	enableLog := false
	if *logFlag == 1 {
		enableLog = true
	}

	// socket rpc时，默认开启日志，方便定位问题
	if *modeFlag == 2 {
		enableLog = true
	}

	// 开启日志时，才开启pprof
	if enableLog {
		go func() {
			http.ListenAndServe("localhost:6060", nil)
		}()
		logger.Init("gop-lsp.log", true)
	}

	// common.GlobalConfigDefautInit()
	// common.GConfig.IntialGlobalVar()

	//*modeFlag = 0
	if *modeFlag == 1 {
		cmdRPC()
	} else if *modeFlag == 2 {
		socketRPC()
	} else if *modeFlag == 0 {
		// runLocalDiagnostices(*localpath)
	}
}

// cmd 的方式运行rpc
func cmdRPC() {
	logger.Debug("local stat running ....")
	Server := langserver.CreateServer()
	Server.Start(channel.Header("")(os.Stdin, os.Stdout))

	logger.Debug("Server started ....")
	if err := Server.Wait(); err != nil {
		logger.Debug("Server exited: %v", err)
	}
	logger.Debug("Server exited return")
}

// 网络的方式运行rpc
func socketRPC() {
	logger.Debug("socket running ....")
	lst, err := net.Listen("tcp", "localhost:7778")
	if err != nil {
		logger.Error(err)
		return
	}

	var wg sync.WaitGroup
	for {
		conn, err := lst.Accept()
		logger.Debug("accept new conn ....")
		Server := langserver.CreateServer()
		if err != nil {
			if channel.IsErrClosing(err) {
				err = nil
			} else {
				logger.Errorf(err, "Error accepting new connection")
			}
			wg.Wait()
			logger.Errorf(err, "Error accepting new connection:")
			return
		}
		ch := channel.Header("")(conn, conn)
		wg.Add(1)
		go func() {
			defer wg.Done()
			Server.Start(ch)
			if err := Server.Wait(); err != nil && err != io.EOF {
				logger.Debugf("Server exited: %v", err)
			}

			logger.Debugf("Server exited11: %v", err)
		}()
	}
}

// func runLocalDiagnostices(localpath string) {
// 	logger.Debug("local Diagnostices running ....")
// 	lspServer := langserver.CreateLspServer()
// 	lspServer.RunLocalDiagnostices(localpath)
// 	logger.Debug("local Diagnostices exited ")
// }
