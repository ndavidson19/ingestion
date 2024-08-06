package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/ndavidson/ingestion/internal/models"
)

type MockProducer struct {
	ProduceChannel chan *kafka.Message
	Events         chan kafka.Event
	ProduceCalled  bool
	CloseCall      int
}

func NewMockProducer() *MockProducer {
	return &MockProducer{
		ProduceChannel: make(chan *kafka.Message, 100),
		Events:         make(chan kafka.Event, 100),
	}
}

func (m *MockProducer) Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error {
	m.ProduceCalled = true
	m.ProduceChannel <- msg
	if deliveryChan != nil {
		go func() {
			deliveryChan <- &kafka.Message{
				TopicPartition: msg.TopicPartition,
				Value:          msg.Value,
				Headers:        msg.Headers,
			}
		}()
	}
	return nil
}

func (m *MockProducer) Events() chan kafka.Event {
	return m.Events
}

func (m *MockProducer) ProduceChannel() chan *kafka.Message {
	return m.ProduceChannel
}

func (m *MockProducer) Close() {
	m.CloseCall++
}

// Implement the Produce method to match our actual Producer interface
func (m *MockProducer) ProduceMsg(data models.IngestData) error {
	m.ProduceCalled = true
	// In a real scenario, we'd convert data to a kafka.Message here
	return nil
}
