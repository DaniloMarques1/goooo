package consumer

import (
	"context"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

type Consumer interface {
	Consume() []byte
}

type KafkaConsumer struct {
	reader *kafka.Reader
}

func NewKafkaConsumer() *KafkaConsumer {
	kConsumer := &KafkaConsumer{}
	kConsumer.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{os.Getenv("KAFKA_BROKER")},
		GroupID:   "register_merchants",
		Topic:     os.Getenv("KAFKA_TOPIC"),
		Partition: 0,
	})
	return kConsumer
}

func (kc *KafkaConsumer) Consume() []byte {
	msg, err := kc.reader.ReadMessage(context.Background())
	if err != nil {
		log.Printf("Error reading from topic %v\n", err)
		return []byte{}
	}
	return msg.Value
}
