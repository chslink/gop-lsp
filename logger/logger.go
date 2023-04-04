package logger

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

var (
	Logger  *log.Logger
	isDebug = false
)

func Init(file string, debug bool) {
	isDebug = debug
	var output io.Writer
	output = os.Stderr
	if file != "" {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		path := filepath.Join(wd, file)
		f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}
		output = f
	}
	Logger = log.New(output, "[gop-lsp]", log.LstdFlags|log.Lshortfile|log.Lmsgprefix)
}

func Printf(format string, v ...any) {
	if Logger != nil {
		Logger.Printf(format, v...)
	}
}

func Println(v ...any) {
	if Logger != nil {
		Logger.Panicln(v...)
	}
}

func Debugf(format string, v ...any) {
	if Logger != nil && isDebug {
		Logger.Printf(format, v...)
	}
}

func Debug(v ...any) {
	if Logger != nil && isDebug {
		Logger.Println(v...)
	}
}

func Infof(format string, v ...any) {
	Printf(format, v...)
}

func Errorf(err error, format string, v ...any) {
	Printf("%+v\n", errors.WithMessagef(err, format, v...))
}

func Error(err error) {
	Printf("%+v\n", errors.WithStack(err))
}
