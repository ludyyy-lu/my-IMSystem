package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

var kafkaWriter *kafka.Writer

func InitKafkaWriter(brokers []string, topic string) {
	kafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func SendMessage(topic string, data []byte) error {
	if kafkaWriter == nil {
		log.Println("Kafka writer is not initialized")
		return nil
	}
	msg := kafka.Message{
		Topic: topic,
		Value: data,
	}
	return kafkaWriter.WriteMessages(context.Background(), msg)
}
