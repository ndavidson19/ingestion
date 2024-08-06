package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type PrometheusMetrics struct {
	requestsProcessed prometheus.Counter
	responseTime      prometheus.Histogram
}

func NewPrometheusMetrics() *PrometheusMetrics {
	return &PrometheusMetrics{
		requestsProcessed: promauto.NewCounter(prometheus.CounterOpts{
			Name: "requests_processed_total",
			Help: "The total number of processed requests",
		}),
		responseTime: promauto.NewHistogram(prometheus.HistogramOpts{
			Name: "response_time_seconds",
			Help: "Response time in seconds",
		}),
	}
}

func (m *PrometheusMetrics) IncrementRequestsProcessed() {
	m.requestsProcessed.Inc()
}

func (m *PrometheusMetrics) ObserveResponseTime(t float64) {
	m.responseTime.Observe(t)
}
