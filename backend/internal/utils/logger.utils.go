package utils

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LogRequestHeaders(c *gin.Context, log *zap.Logger) {
	fields := []zap.Field{
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
	}

	if userID, exists := c.Get("user_id"); exists {
		fields = append(fields, zap.String("user_id", userID.(string)))
	}
	if role, exists := c.Get("role"); exists {
		fields = append(fields, zap.String("role", role.(string)))
	}

	if requestID := c.GetHeader("X-Request-ID"); requestID != "" {
		fields = append(fields, zap.String("request_id", requestID))
	}

	log.Info("request", fields...)
}
