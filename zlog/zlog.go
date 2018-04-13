package zlog

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// SetColorZapGlobalLogger -
func SetColorZapGlobalLogger() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zlog, err := config.Build()
	if err != nil {
		log.Fatal("set color zap global logger: err - ", err)
	}
	zap.ReplaceGlobals(zlog)
}
