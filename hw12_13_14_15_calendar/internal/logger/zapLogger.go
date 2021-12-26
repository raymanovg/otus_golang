package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLogger(output []string, level string, devMode bool) (*zap.Logger, error) {
	atomicLevel := zap.NewAtomicLevel()
	err := atomicLevel.UnmarshalText([]byte(level))
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal log level string: %w", err)
	}

	conf := zap.Config{
		Level:       atomicLevel,
		Encoding:    "json",
		Development: devMode,
		OutputPaths: output,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "level",
			TimeKey:        "timestamp",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
		},
	}

	return conf.Build()
}
