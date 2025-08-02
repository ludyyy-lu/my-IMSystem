package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/segmentio/kafka-go"
)

var kafkaWriter *kafka.Writer

func InitKafkaProducer(brokers []string) {
	kafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Balancer: &kafka.LeastBytes{},
		Async:    true,
	}
}

func SendMessage(topic string, value any) error {
	if kafkaWriter == nil {
		log.Println("Kafka writer is not initialized")
		return errors.New("kafka writer not initialized")
	}

	// 序列化 value
	msgBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Topic: topic,
		Value: msgBytes,
	}
	return kafkaWriter.WriteMessages(context.Background(), msg)
}
