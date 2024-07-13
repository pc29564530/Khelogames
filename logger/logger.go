package logger

import (
	"os"

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
	})
	log.SetLevel(logrus.InfoLevel)
	log.SetLevel(logrus.ErrorLevel)
	log.SetLevel(logrus.DebugLevel)

	return &Logger{log}
}
