package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ndavidson/ingestion/internal/handlers"
	"github.com/ndavidson/ingestion/internal/kafka"
	"github.com/ndavidson/ingestion/internal/models"
	"github.com/ndavidson/ingestion/pkg/metrics"
)

func TestIngestHandler(t *testing.T) {
	// Mock Kafka producer
	mockProducer := &kafka.MockProducer{}

	// Initialize metrics
	metrics := metrics.NewPrometheusMetrics()

	// Create a request to pass to our handler
	ingestData := models.IngestData{
		Source: "test",
		Data: map[string]interface{}{
			"key": "value",
		},
	}
	body, _ := json.Marshal(ingestData)
	req, err := http.NewRequest("POST", "/v1/ingest", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.IngestHandler(mockProducer, metrics))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusAccepted)
	}

	// Check the response body is what we expect.
	expected := `{"status":"accepted"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Check that the mock producer was called
	if !mockProducer.ProduceCalled {
		t.Errorf("Expected Produce method to be called")
	}
}
