package logger

import (
	"fmt"

	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLogger(conf config.LoggerConf) (*zap.Logger, error) {
	atomicLevel := zap.NewAtomicLevel()
	err := atomicLevel.UnmarshalText([]byte(conf.Level))
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal log level string: %w", err)
	}

	zapConf := zap.Config{
		Level:       atomicLevel,
		Encoding:    "json",
		Development: conf.DevMode,
		OutputPaths: conf.Output,
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

	logger, err := zapConf.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger from config: %w", err)
	}

	return logger.With(zap.String("version", config.VERSION)), nil
}
