// File: cmd/rest-ingestion-service/main.go

package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ndavidson/ingestion/internal/config"
	"github.com/ndavidson/ingestion/internal/handlers"
	"github.com/ndavidson/ingestion/internal/kafka"
	"github.com/ndavidson/ingestion/internal/middleware"
	"github.com/ndavidson/ingestion/pkg/logger"
	"github.com/ndavidson/ingestion/pkg/metrics"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Initialize logger
	logger := logger.NewLogger()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration", "error", err)
	}

	// Initialize Kafka producer
	producer, err := kafka.NewProducer(cfg.KafkaBrokers)
	if err != nil {
		logger.Fatal("Failed to create Kafka producer", "error", err)
	}
	defer producer.Close()

	// Initialize metrics
	metrics := metrics.NewPrometheusMetrics()

	// Set up HTTP server
	r := mux.NewRouter()

	// Middleware
	r.Use(middleware.Logging(logger))
	r.Use(middleware.RateLimit(cfg.RateLimit))
	r.Use(middleware.Auth(cfg.AuthToken))

	// Routes
	r.HandleFunc("/v1/ingest", handlers.IngestHandler(producer, metrics)).Methods("POST")
	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/health", handlers.HealthCheckHandler).Methods("GET")

	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		logger.Info("Starting server", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", "error", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	// Implement graceful shutdown logic here
}
