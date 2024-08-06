package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ndavidson/ingestion/internal/kafka"
	"github.com/ndavidson/ingestion/internal/models"
	"github.com/ndavidson/ingestion/pkg/metrics"
)

func TestIngestHandlerSuccess(t *testing.T) {
	mockProducer := kafka.MockProducer()
	metrics := metrics.NewPrometheusMetrics()

	ingestData := models.IngestData{
		Source: "test",
		Data: map[string]interface{}{
			"key": "value",
		},
	}
	body, _ := json.Marshal(ingestData)
	req := httptest.NewRequest("POST", "/v1/ingest", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := IngestHandler(mockProducer, metrics)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusAccepted)
	}

	expected := `{"status":"accepted"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	if !mockProducer.ProduceCalled {
		t.Errorf("Expected Produce method to be called")
	}
}

func TestIngestHandlerInvalidJSON(t *testing.T) {
	mockProducer := &kafka.MockProducer{}
	metrics := metrics.NewPrometheusMetrics()

	req := httptest.NewRequest("POST", "/v1/ingest", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := IngestHandler(mockProducer, metrics)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}
