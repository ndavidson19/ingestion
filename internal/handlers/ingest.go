package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ndavidson/ingestion/internal/kafka"
	"github.com/ndavidson/ingestion/internal/models"
	"github.com/ndavidson/ingestion/pkg/metrics"
)

func IngestHandler(producer *kafka.Producer, metrics *metrics.PrometheusMetrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		var ingestData models.IngestData
		if err := json.NewDecoder(r.Body).Decode(&ingestData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Validate and process ingestData...

		if err := producer.Produce(ingestData); err != nil {
			http.Error(w, "Failed to publish message", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{"status": "accepted"})

		metrics.IncrementRequestsProcessed()
		metrics.ObserveResponseTime(time.Since(startTime).Seconds())
	}
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}
