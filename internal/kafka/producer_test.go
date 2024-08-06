// File: internal/kafka/mock_producer.go

package kafka

import (
	"testing"

	"github.com/ndavidson/ingestion/internal/models"
)

type MockProducer struct {
	ProduceCalled bool
	ProduceErr    error
}

func (m *MockProducer) Produce(data models.IngestData) error {
	m.ProduceCalled = true
	return m.ProduceErr
}

func (m *MockProducer) Close() {}

func TestProducer_Produce(t *testing.T) {
	mockProducer := &MockProducer{}
	producer := &Producer{producer: mockProducer}

	data := models.IngestData{
		Source: "test",
		Data:   map[string]interface{}{"key": "value"},
	}

	err := producer.Produce(data)
	if err != nil {
		t.Fatalf("Produce returned an error: %v", err)
	}

	if !mockProducer.ProduceCalled {
		t.Error("Expected Produce to be called on the underlying Kafka producer")
	}
}
