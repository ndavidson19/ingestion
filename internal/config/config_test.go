package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("KAFKA_BROKERS", "localhost:9092")
	os.Setenv("PORT", "8080")
	os.Setenv("AUTH_TOKEN", "test-token")
	os.Setenv("RATE_LIMIT", "100")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.KafkaBrokers != "localhost:9092" {
		t.Errorf("Expected KafkaBrokers to be 'localhost:9092', got '%s'", cfg.KafkaBrokers)
	}
	if cfg.Port != "8080" {
		t.Errorf("Expected Port to be '8080', got '%s'", cfg.Port)
	}
	if cfg.AuthToken != "test-token" {
		t.Errorf("Expected AuthToken to be 'test-token', got '%s'", cfg.AuthToken)
	}
	if cfg.RateLimit != 100 {
		t.Errorf("Expected RateLimit to be 100, got %d", cfg.RateLimit)
	}
}
