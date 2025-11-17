package logger

import (
	"os"
	"time"
	"webserver/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

func Init(cfg *config.Config) error {
	var level zapcore.Level
	switch cfg.LogLevel {
	case "debug":
		level = zap.DebugLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		level,
	)

	globalLogger = zap.New(core, zap.AddCaller())
	return nil
}

func Get() *zap.Logger {
	if globalLogger == nil {
		// Fallback к базовому логгеру если не инициализирован
		basicLogger, _ := zap.NewProduction()
		return basicLogger
	}
	return globalLogger
}

func Sync() {
	if globalLogger != nil {
		_ = globalLogger.Sync()
	}
}

// Оберточные функции для удобства
func String(key string, val string) zap.Field {
	return zap.String(key, val)
}

func Error(err error) zap.Field {
	return zap.Error(err)
}

func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

func Int64(key string, val int64) zap.Field {
	return zap.Int64(key, val)
}

func Duration(key string, val time.Duration) zap.Field {
	return zap.Duration(key, val)
}

func Any(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}
