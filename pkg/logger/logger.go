package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type contextKey struct{}

func InitLogger(env string) (*zap.Logger, zap.AtomicLevel, error) {
	level := zap.NewAtomicLevelAt(zap.InfoLevel)

	if env == "development" {
		level.SetLevel(zap.DebugLevel)
	}

	config := zap.NewProductionConfig()
	config.Level = level
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build(zap.AddCaller(),zap.AddCallerSkip(1))
	return logger, level, err
}

// WithContext stores the logger inside a standard context.Context
func WithContext(ctx context.Context, log *zap.Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, log)
}

// FromContext retrieves the logger from context.
// Falls back to a no-op logger so callers never need to nil-check.
func FromContext(ctx context.Context) *zap.Logger {
	if log, ok := ctx.Value(contextKey{}).(*zap.Logger); ok && log != nil {
		return log
	}
	return zap.NewNop()
}