package middleware

import (
	"payment-service/internal/observability"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		duration := time.Since(start).Seconds()

		status := strconv.Itoa(c.Writer.Status())
		path := c.FullPath()

		observability.HTTPRequests.WithLabelValues(
			c.Request.Method,
			path,
			status,
		).Inc()

		observability.HTTPDuration.WithLabelValues(
			c.Request.Method,
			path,
			status,
		).Observe(duration)

	}
}
