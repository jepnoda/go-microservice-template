package logger

import (
	"sync"

	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	once   sync.Once
)

// GetLogger returns a singleton zap.Logger instance configured for console output
func GetLogger() *zap.Logger {
	once.Do(func() {
		cfg := zap.NewDevelopmentConfig()
		cfg.Encoding = "console"
		var err error
		logger, err = cfg.Build()
		if err != nil {
			panic(err)
		}
		logger.Info("Logger initialized")
	})
	return logger
}

func Info(msg string, fields ...zap.Field) {
	if GetLogger() != nil {
		GetLogger().Info(msg, fields...)
	}
}

func Error(msg string, fields ...zap.Field) {
	if GetLogger() != nil {
		GetLogger().Error(msg, fields...)
	}
}
