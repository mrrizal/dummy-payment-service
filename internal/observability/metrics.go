package observability

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	HTTPRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status", "payment_method"},
	)

	HTTPDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status", "payment_method"},
	)

	DBQueryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Database query latency",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	DBErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_errors_total",
			Help: "Total database errors",
		},
		[]string{"operation"},
	)
)

func InitMetrics() {
	prometheus.MustRegister(HTTPRequests)
	prometheus.MustRegister(HTTPDuration)
	prometheus.MustRegister(DBQueryDuration)
	prometheus.MustRegister(DBErrors)
}
