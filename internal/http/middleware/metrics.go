package middleware

import (
	"payment-service/internal/observability"
	"strconv"

	"github.com/gin-gonic/gin"
)

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		path := c.FullPath()

		observability.HTTPRequests.WithLabelValues(
			c.Request.Method,
			path,
			status,
		).Inc()
	}
}
