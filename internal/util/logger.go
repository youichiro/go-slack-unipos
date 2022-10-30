package util

import "go.uber.org/zap"

var Logger *zap.Logger

func SetLogger() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	Logger = logger
	return logger
}
