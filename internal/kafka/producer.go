package kafka

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/ndavidson/ingestion/internal/models"
)

type Producer struct {
	producer *kafka.Producer
}

func NewProducer(brokers string) (*Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": brokers})
	if err != nil {
		return nil, err
	}
	return &Producer{producer: p}, nil
}

func (p *Producer) Produce(data models.IngestData) error {
	topic := "raw_data_" + data.Source
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value,
		Timestamp:      data.Timestamp,
	}, nil)
}

func (p *Producer) Close() {
	p.producer.Close()
}
