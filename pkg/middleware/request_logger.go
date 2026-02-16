package middleware

import (
	"github.com/adarsh-jaiss/the-bridge/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const loggerKey = "logger"

func RequestLogger(baseLogger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use existing request ID if present (from proxy)
		requestId := c.GetHeader("X-Request-ID")
		if requestId == "" {
			requestId = uuid.New().String()
		}

		// Attach structured fields
		reqLogger := baseLogger.With(
			zap.String("request_id", requestId),
			zap.String("path", c.FullPath()),
			zap.String("method", c.Request.Method),
			zap.String("client_ip", c.ClientIP()),
		)

		//store in context
		// Inject into standard context so it flows into every layer
		ctx := logger.WithContext(c.Request.Context(), reqLogger)
		c.Request = c.Request.WithContext(ctx)

		// Return request id to client
		c.Writer.Header().Set("X-Request-ID", requestId)
		c.Next()
	}
}

func GetLogger(c *gin.Context) *zap.Logger {
	if logger, exists := c.Get(loggerKey); exists {
		return logger.(*zap.Logger)
	}
	return zap.NewNop()
}