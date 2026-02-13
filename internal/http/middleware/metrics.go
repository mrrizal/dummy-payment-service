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

		// default value
		paymentMethod := "n/a"

		// ambil dari context kalau ada
		if pm, exists := c.Get("payment_method"); exists {
			if val, ok := pm.(string); ok {
				paymentMethod = val
			}
		}

		observability.HTTPRequests.WithLabelValues(
			c.Request.Method,
			path,
			status,
			paymentMethod,
		).Inc()

		observability.HTTPDuration.WithLabelValues(
			c.Request.Method,
			path,
			status,
			paymentMethod,
		).Observe(duration)

	}
}
