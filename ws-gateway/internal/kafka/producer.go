package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

var producer *kafka.Writer

func InitKafkaProducer(brokers []string, topic string) {
	producer = kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
		Async:   true, // 异步发送
	})
}

func SendMessage(value interface{}) error {
	msgBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = producer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   nil,
			Value: msgBytes,
		},
	)
	if err != nil {
		log.Printf("Failed to send Kafka message: %v", err)
	}
	return err
}
