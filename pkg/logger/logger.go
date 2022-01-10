package logger

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
)

//go:generate mockgen -source=logger.go -destination=mocks/mock_logger.go -package=logger_mock

// Logger -.
type Logger struct {
	logger *logrus.Logger
}

// Interface -.
type Interface interface {
	Debug(message interface{}, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
}

func New(logLevel string) *Logger {

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
	}
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
		panic("Not supported")
	}
}

// Debug -.
func (l *Logger) Debug(message interface{}, args ...interface{}) {
	l.msg("debug", message, args...)
}

// Info -.
func (l *Logger) Info(message string, args ...interface{}) {
	l.log(message, args...)
}

// Warn -.
func (l *Logger) Warn(message string, args ...interface{}) {
	l.log(message, args...)
}

// Error -.
func (l *Logger) Error(message interface{}, args ...interface{}) {
	if l.logger.GetLevel() == logrus.DebugLevel {
		l.Debug(message, args...)
	} else {
		l.msg("error", message, args...)
	}
}

// Fatal -.
func (l *Logger) Fatal(message interface{}, args ...interface{}) {
	l.msg("fatal", message, args...)

	os.Exit(1)
}

func (l *Logger) log(message string, args ...interface{}) {
	if len(args) == 0 {
		l.logger.Info(message)
	} else {
		l.logger.Infof(message, args...)
	}
}
func (l *Logger) msg(level string, message interface{}, args ...interface{}) {
	switch msg := message.(type) {
	case error:
		l.log(msg.Error(), args...)
	case string:
		l.log(msg, args...)
	default:
		l.log(fmt.Sprintf("%s message %v has unknown type %v", level, message, msg), args...)
	}
}
