package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Field = zap.Field

type Logger struct {
	*zap.Logger
}

func New(level string) *Logger {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(parseLevel(level))
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	log, _ := cfg.Build(zap.AddCallerSkip(1))
	return &Logger{Logger: log}
}

func parseLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// 便捷方法
func String(key, value string) Field      { return zap.String(key, value) }
func Int(key string, value int) Field     { return zap.Int(key, value) }
func Int64(key string, value int64) Field { return zap.Int64(key, value) }
func Float64(key string, value float64) Field { return zap.Float64(key, value) }
func Bool(key string, value bool) Field   { return zap.Bool(key, value) }
func Any(key string, value interface{}) Field { return zap.Any(key, value) }
func Error(err error) Field               { return zap.Error(err) }

func Duration(key string, value interface{}) Field {
	switch v := value.(type) {
	case int64:
		return zap.Int64(key, v)
	default:
		return zap.Any(key, value)
	}
}
