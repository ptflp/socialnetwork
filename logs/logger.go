package logs

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	GeneralError = 2
)

func NewLogger() *zap.Logger {
	zapConf := zap.NewProductionConfig()
	zapConf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapConf.EncoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	atom := zap.NewAtomicLevelAt(zapcore.DebugLevel)
	zapConf.Level = atom
	logger, err := zapConf.Build()
	if err != nil {
		fmt.Printf("logs initialization: %s\n", err)
		os.Exit(GeneralError)
	}

	return logger.Named("infoBlog")
}
