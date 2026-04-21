package configs

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func InitLogger() *zap.Logger {
	var config zap.Config

	if os.Getenv("GIN_MODE") == "release" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	Log = logger
	return Log
}

func GetLogger() *zap.Logger {
	if Log == nil {
		return InitLogger()
	}
	return Log
}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
