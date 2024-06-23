package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger, _ = zap.Config{
	Encoding:    "json",
	Level:       zap.NewAtomicLevelAt(zapcore.InfoLevel),
	OutputPaths: []string{"stdout"},
	EncoderConfig: zapcore.EncoderConfig{
		MessageKey:  "message",
		EncodeLevel: zapcore.CapitalLevelEncoder, // INFO

		TimeKey:    "time",
		EncodeTime: zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
	},
}.Build()
