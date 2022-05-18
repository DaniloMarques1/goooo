package consumer

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader *kafka.Reader
}

func NewKafkaConsumer() *KafkaConsumer {
	kConsumer := &KafkaConsumer{}
	kConsumer.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"}, // @@@
		GroupID:   "registerMerchants",
		Topic:     "merchants", // @@@
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
