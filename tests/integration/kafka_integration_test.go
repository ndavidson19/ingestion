package integration

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/ndavidson/ingestion/internal/models"
)

func TestKafkaIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		t.Fatalf("Failed to create producer: %s", err)
	}
	defer producer.Close()

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "test-group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		t.Fatalf("Failed to create consumer: %s", err)
	}
	defer consumer.Close()

	topic := "test_topic"
	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		t.Fatalf("Failed to subscribe to topic: %s", err)
	}

	testData := models.IngestData{
		Source:    "test",
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"key": "value",
		},
	}

	value, _ := json.Marshal(testData)
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value,
	}, nil)
	if err != nil {
		t.Fatalf("Failed to produce message: %s", err)
	}

	msg, err := consumer.ReadMessage(10 * time.Second)
	if err != nil {
		t.Fatalf("Failed to consume message: %s", err)
	}

	var receivedData models.IngestData
	err = json.Unmarshal(msg.Value, &receivedData)
	if err != nil {
		t.Fatalf("Failed to unmarshal message: %s", err)
	}

	if receivedData.Source != testData.Source {
		t.Errorf("Received data does not match sent data")
	}
}
