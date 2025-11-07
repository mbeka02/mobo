package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// type contextKey string

// const requestIDKey contextKey = "request_id"

var globalLogger *zap.Logger

func Init(env string) error {
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logger, err := config.Build(
		zap.AddCallerSkip(1), // Skip wrapper functions in stack trace
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return err
	}

	globalLogger = logger
	return nil
}

// Get returns the global logger instance
func Get() *zap.Logger {
	if globalLogger == nil {
		// Fallback logger if init wasn't called
		globalLogger, _ = zap.NewProduction()
	}
	return globalLogger
}

// Sync flushes any buffered log entries (use during shutdown)
func Sync() error {
	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
}

// WithRequestID adds request ID to context for tracing
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// FromContext returns a logger with context fields (e.g., request_id)
func FromContext(ctx context.Context) *zap.Logger {
	logger := Get()

	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		logger = logger.With(zap.String("request_id", requestID))
	}

	return logger
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	Get().Info(msg, fields...)
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	Get().Error(msg, fields...)
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	Get().Warn(msg, fields...)
}

// Debug logs a debug message
func Debug(msg string, fields ...zap.Field) {
	Get().Debug(msg, fields...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...zap.Field) {
	Get().Fatal(msg, fields...)
}

// InfoCtx logs with context
func InfoCtx(ctx context.Context, msg string, fields ...zap.Field) {
	FromContext(ctx).Info(msg, fields...)
}

// ErrorCtx logs error with context
func ErrorCtx(ctx context.Context, msg string, fields ...zap.Field) {
	FromContext(ctx).Error(msg, fields...)
}

// WarnCtx logs warning with context
func WarnCtx(ctx context.Context, msg string, fields ...zap.Field) {
	FromContext(ctx).Warn(msg, fields...)
}

// DebugCtx logs debug with context
func DebugCtx(ctx context.Context, msg string, fields ...zap.Field) {
	FromContext(ctx).Debug(msg, fields...)
}
