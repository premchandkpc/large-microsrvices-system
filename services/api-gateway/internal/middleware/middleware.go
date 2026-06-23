package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		requestID, _ := c.Get("request_id")

		fields := []zap.Field{
			zap.Int("status", status),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Duration("latency", latency),
			zap.String("request_id", requestID.(string)),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("client_ip", c.ClientIP()),
		}

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				logger.Error("request error", append(fields, zap.Error(e.Err))...)
			}
		} else if status >= 500 {
			logger.Error("server error", fields...)
		} else if status >= 400 {
			logger.Warn("client error", fields...)
		} else {
			logger.Info("request", fields...)
		}
	}
}

func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID, _ := c.Get("request_id")
				logger.Error("panic recovered",
					zap.Any("error", err),
					zap.String("request_id", requestID.(string)),
					zap.Stack("stack"),
				)
				c.AbortWithStatusJSON(500, gin.H{
					"error":      "internal server error",
					"request_id": requestID,
				})
			}
		}()
		c.Next()
	}
}
