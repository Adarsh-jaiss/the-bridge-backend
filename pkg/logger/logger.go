package logger

import (
	"context"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type contextKey struct{}

func InitLogger(env string) (*zap.Logger, zap.AtomicLevel, error) {
	level := zap.NewAtomicLevelAt(zap.InfoLevel)

	var config zap.Config

	if env == "development" {
		config = zap.NewDevelopmentConfig()
		level.SetLevel(zap.DebugLevel)

		// Pretty logs
		config.Encoding = "console"
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config = zap.NewProductionConfig()
		config.Encoding = "json"
	}

	// Clean timestamp
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}

	// Clean caller
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	//  Disable default stacktrace
	config.DisableStacktrace = true

	logger, err := config.Build(
		zap.AddCaller(),
		zap.AddStacktrace(zap.PanicLevel))
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
