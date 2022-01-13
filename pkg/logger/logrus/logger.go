package logger

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
)

type Logger struct {
	logger *logrus.Logger
}

func New(logLevel string) (*Logger, error) {

	var level logrus.Level
	level = LogLevel(logLevel)
	logger := &logrus.Logger{
		Out:   os.Stdout,
		Level: level,
		Formatter: &prefixed.TextFormatter{
			DisableColors:   true,
			TimestampFormat: "2009-06-03 11:04:075",
		},
	}
	return &Logger{
		logger: logger,
	}, nil
}

func LogLevel(lvl string) logrus.Level {
	switch lvl {
	case "info":
		return logrus.InfoLevel
	case "error":
		return logrus.ErrorLevel
	case "debug":
		return logrus.DebugLevel
	default:
		return logrus.InfoLevel
	}
}

func (l *Logger) Debug(message string, args ...interface{}) {
	l.log(message, args...)
}

func (l *Logger) Info(message string, args ...interface{}) {
	l.log(message, args...)
}

func (l *Logger) Warn(message string, args ...interface{}) {
	l.log(message, args...)
}

func (l *Logger) Error(message string, args ...interface{}) {
	if l.logger.GetLevel() == logrus.DebugLevel {
		l.Debug(message, args...)
	} else {
		l.log(message, args...)
	}
}

func (l *Logger) Fatal(message string, args ...interface{}) {
	l.log(message, args...)

	os.Exit(1)
}

func (l *Logger) log(message string, args ...interface{}) {
	if len(args) == 0 {
		l.logger.Info(message)
	} else {
		l.logger.Infof(message, args...)
	}
}
