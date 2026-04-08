package kafka

import (
"context"
"encoding/json"
"fmt"
"time"

"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
writer *kafka.Writer
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
return &KafkaProducer{
writer: &kafka.Writer{
Addr:         kafka.TCP(brokers...),
Topic:        topic,
Balancer:     &kafka.LeastBytes{},
RequiredAcks: kafka.RequireAll,
Async:        false,
},
}
}

// SendMessage serialises value as JSON (or sends raw bytes if already []byte)
// and writes it to Kafka with the given key.
func (p *KafkaProducer) SendMessage(key string, value interface{}) error {
var payload []byte
switch v := value.(type) {
case []byte:
payload = v
default:
var err error
payload, err = json.Marshal(v)
if err != nil {
return fmt.Errorf("kafka producer: marshal failed: %w", err)
}
}
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
return p.writer.WriteMessages(ctx, kafka.Message{
Key:   []byte(key),
Value: payload,
})
}
