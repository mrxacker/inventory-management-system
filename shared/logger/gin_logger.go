// shared/logger/gin_logger.go

package logger

import (
	"time"

	"github.com/gin-gonic/gin"
)

func GinLogger(logger Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get status code
		statusCode := c.Writer.Status()

		// Get client IP
		clientIP := c.ClientIP()

		// Get request method
		method := c.Request.Method

		// Get user agent
		userAgent := c.Request.UserAgent()

		// Get errors if any
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		fields := []interface{}{
			"status", statusCode,
			"method", method,
			"path", path,
			"query", query,
			"ip", clientIP,
			"user_agent", userAgent,
			"latency", latency.Milliseconds(),
			"latency_human", latency.String(),
		}

		if errorMessage != "" {
			fields = append(fields, "error", errorMessage)
		}

		// Log based on status code
		switch {
		case statusCode >= 500:
			logger.Error("Server error", fields...)
		case statusCode >= 400:
			logger.Warn("Client error", fields...)
		case statusCode >= 300:
			logger.Info("Redirection", fields...)
		default:
			logger.Info("Success", fields...)
		}
	}
}

// GinRecovery returns a gin middleware for recovering from panics
func GinRecovery(logger Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered",
					"error", err,
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
				)
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}
