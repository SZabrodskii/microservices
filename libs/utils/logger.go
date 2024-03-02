package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Logger struct {
	logger *zap.Logger
}

type LoggerConfig struct {
	Level       string   `json:"level"`
	Encoding    string   `json:"encoding"`
	OutputPaths []string `json:"output paths"`
	Development bool     `json:"development"`
}

func NewLogger(config LoggerConfig) (*Logger, error) {
	var lvl zapcore.Level
	err := lvl.UnmarshalText([]byte(config.Level))
	if err != nil {
		return nil, err
	}
	encoderConfig := zap.NewProductionEncoderConfig()
	if config.Development {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	}

	var encoder zapcore.Encoder
	switch config.Encoding {
	case "json":
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	//case "bson":
	//	encoder = zapcore.NewConsoleEncoder(encoderConfig)
	default:
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	var cores []zapcore.Core
	for _, outputPath := range config.OutputPaths {
		output, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			output = os.Stderr
		}
		core := zapcore.NewCore(encoder, zapcore.AddSync(output), lvl)
		cores = append(cores, core)
	}

	multiCore := zapcore.NewTee(cores...)
	loggerInstance := zap.New(multiCore)

	return &Logger{
		logger: loggerInstance,
	}, nil
}
