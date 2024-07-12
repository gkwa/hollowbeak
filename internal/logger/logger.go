package logger

import (
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewConsoleLogger(verbosity int, json bool) logr.Logger {
	var zapLogger *zap.Logger
	var err error

	config := zap.NewProductionConfig()
	if !json {
		config.Encoding = "console"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	config.Level = zap.NewAtomicLevelAt(zapcore.Level(-verbosity))

	zapLogger, err = config.Build()
	if err != nil {
		panic(err)
	}

	return zapr.NewLogger(zapLogger)
}
