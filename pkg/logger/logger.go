package logger

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger    *zap.Logger
	startTime time.Time

	initOnce sync.Once
)

// ZapLoggerInit starts logger instance and configures it. Reads env variable IS_PROD and based on its presence sets log level to print
func ZapLoggerInit() {
	initOnce.Do(func() {
		startTime = time.Now()
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder := zapcore.NewJSONEncoder(encoderConfig)
		writer := zapcore.AddSync(os.Stderr)

		//prod - only error+
		if strings.Contains(os.Getenv("IS_PROD"), "true") {

			logLevel := zapcore.ErrorLevel
			core := zapcore.NewTee(
				zapcore.NewCore(encoder, writer, logLevel),
			)

			logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))
			return
		}

		//dev mode debug+
		logLevel := zapcore.DebugLevel
		core := zapcore.NewTee(
			zapcore.NewCore(encoder, writer, logLevel),
		)

		logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	})
}

func Info(message string, fields ...zap.Field) {
	logger.Info(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	logger.Debug(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
	logger.Warn(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	logger.Error(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	logger.Fatal(message, fields...)
}

func ServiceError(custom, err error, data ...any) {
	Error(fmt.Sprintf("%v, err - %v, info - %+v", custom, err, data))
}

func Uptime() time.Duration {
	return time.Since(startTime)
}
