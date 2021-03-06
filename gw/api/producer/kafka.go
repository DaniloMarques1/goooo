package producer

import (
	"context"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	conn *kafka.Conn
}

func NewKafkaProducer() (*KafkaProducer, error) {
	topic := os.Getenv("TOPIC")
	addr := os.Getenv("KAFKA_ADDR")
	partition := 0
	conn, err := kafka.DialLeader(context.Background(), "tcp", addr, topic, partition)
	if err != nil {
		return nil, err
	}
	return &KafkaProducer{conn: conn}, nil
}

func (kp *KafkaProducer) Produce(b []byte) error {
	_, err := kp.conn.WriteMessages(
		kafka.Message{
			Value: b,
		},
	)

	if err != nil {
		log.Printf("Error producing message %v\n", err)
		return err
	}
	return nil
}

func (kp *KafkaProducer) Close() error {
	if err := kp.conn.Close(); err != nil {
		log.Printf("Error closing kafka connection %v\n", err)
		return err
	}
	return nil
}
