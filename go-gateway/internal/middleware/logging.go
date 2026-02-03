package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go-gateway/pkg/logger"
)

func Logging(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		cost := time.Since(start)
		status := c.Writer.Status()

		fields := []logger.Field{
			logger.String("method", c.Request.Method),
			logger.String("path", path),
			logger.String("query", query),
			logger.Int("status", status),
			logger.String("latency", cost.String()),
			logger.String("client_ip", c.ClientIP()),
		}

		if requestID, exists := c.Get(ContextRequestIDKey); exists {
			fields = append(fields, logger.String("request_id", requestID.(string)))
		}

		if status >= 500 {
			log.Error("request", fields...)
		} else if status >= 400 {
			log.Warn("request", fields...)
		} else {
			log.Info("request", fields...)
		}
	}
}
