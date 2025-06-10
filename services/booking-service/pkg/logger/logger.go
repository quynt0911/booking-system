package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type LoggerInterface interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
	Warn(args ...interface{})
	Fatal(args ...interface{})
	WithFields(fields logrus.Fields) *logrus.Entry
}

type loggerImpl struct {
	*logrus.Logger
}

func NewLogger() LoggerInterface {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)
	return &loggerImpl{logger}
}

func (l *loggerImpl) Debug(args ...interface{}) {
	l.Logger.Debug(args...)
}

func (l *loggerImpl) Info(args ...interface{}) {
	l.Logger.Info(args...)
}

func (l *loggerImpl) Error(args ...interface{}) {
	l.Logger.Error(args...)
}

func (l *loggerImpl) Warn(args ...interface{}) {
	l.Logger.Warn(args...)
}

func (l *loggerImpl) Fatal(args ...interface{}) {
	l.Logger.Fatal(args...)
}

func (l *loggerImpl) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.Logger.WithFields(fields)
}
