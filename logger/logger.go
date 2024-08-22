package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func NewLogger() *Logger {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-01 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			filename := filepath.Base(f.File)
			return filename, fmt.Sprintf(":%d", f.Line)
		},
	})
	log.SetReportCaller(true)
	log.SetLevel(logrus.InfoLevel)
	log.SetLevel(logrus.ErrorLevel)
	log.SetLevel(logrus.DebugLevel)

	return &Logger{log}
}
