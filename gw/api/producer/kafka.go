package producer

import (
	"context"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	//writer *kafka.Writer
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
	//writer := &kafka.Writer{
	//	Addr:  kafka.TCP(addr),
	//	Topic: topic,
	//}

	//return &KafkaProducer{writer: writer}
	return &KafkaProducer{conn: conn}, nil
}

func (kp *KafkaProducer) Produce(merchant Merchant) error {
	b, err := merchant.Marshal()
	if err != nil {
		return err
	}
	_, err = kp.conn.WriteMessages(
		kafka.Message{
			Key:   []byte("merchant"),
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
