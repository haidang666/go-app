package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger      *zap.Logger
	sugar       *zap.SugaredLogger
	initOnce    sync.Once
	loggerError error
)

func initLogger() {
	env := os.Getenv("APP_ENV")

	var cfg zap.Config
	if env == "production" {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "timestamp"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}

	logger, loggerError = cfg.Build()
	if loggerError != nil {
		panic("failed to initialize logger: " + loggerError.Error())
	}
	sugar = logger.Sugar().WithOptions(zap.AddStacktrace(zap.DPanicLevel))
}

func L() *zap.SugaredLogger {
	initOnce.Do(initLogger)
	return sugar
}
