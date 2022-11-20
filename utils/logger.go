package utils

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger(logPath string) *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "ts"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	atomicLevel := zap.NewAtomicLevelAt(zap.InfoLevel)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(&lumberjack.Logger{
			Filename:   fmt.Sprintf("%s/%s.log", logPath, "data"),
			MaxSize:    10, // megabytes
			MaxBackups: 10,
			MaxAge:     30, // days
		}),
		atomicLevel,
	)
	return zap.New(core, zap.AddCaller())
}

func NewRawLogger(logPath string) *zap.Logger {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	atomicLevel := zap.NewAtomicLevelAt(zap.InfoLevel)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(&lumberjack.Logger{
			Filename: fmt.Sprintf("%s/%s.log", logPath, "data"),
			MaxSize:  20, // megabytes
		}),
		atomicLevel,
	)
	return zap.New(core, zap.AddCaller())
}
