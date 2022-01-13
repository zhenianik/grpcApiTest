package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.Logger
}

func New(logLevel string) (*Logger, error) {

	var level zapcore.Level
	level = LogLevel(logLevel)

	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(level),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{
		logger: logger,
	}, nil
}

func LogLevel(lvl string) zapcore.Level {
	switch lvl {
	case "info":
		return zapcore.InfoLevel
	case "error":
		return zapcore.ErrorLevel
	case "debug":
		return zapcore.DebugLevel
	default:
		return zapcore.InfoLevel
	}
}

func (l *Logger) Debug(message string, args ...interface{}) {
	l.logger.Debug(message)
}

func (l *Logger) Info(message string, args ...interface{}) {
	l.logger.Info(message)
}

func (l *Logger) Warn(message string, args ...interface{}) {
	l.logger.Warn(message)
}

func (l *Logger) Error(message string, args ...interface{}) {
	l.logger.Error(message)
}

func (l *Logger) Fatal(message string, args ...interface{}) {
	l.logger.Fatal(message)
}
