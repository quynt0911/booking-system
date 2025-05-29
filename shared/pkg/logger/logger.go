package logger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger(env string) {
	Log = logrus.New()

	// Set log format
	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	// Set output
	Log.SetOutput(os.Stdout)

	// Set log level based on environment
	switch env {
	case "production":
		Log.SetLevel(logrus.WarnLevel)
	case "staging":
		Log.SetLevel(logrus.InfoLevel)
	default:
		Log.SetLevel(logrus.DebugLevel)
	}
}

func Info(args ...interface{}) {
	Log.Info(args...)
}

func Debug(args ...interface{}) {
	Log.Debug(args...)
}

func Error(args ...interface{}) {
	Log.Error(args...)
}

func Warn(args ...interface{}) {
	Log.Warn(args...)
}

func Fatal(args ...interface{}) {
	Log.Fatal(args...)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return Log.WithFields(fields)
}
