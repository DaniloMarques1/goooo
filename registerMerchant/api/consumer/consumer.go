package consumer

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/danilomarques1/godemo/registerMerchant/api/model"
	"github.com/segmentio/kafka-go"
)

type Consumer interface {
	Consume() []byte
}

type KafkaConsumer struct {
	reader             *kafka.Reader
	merchantRepository model.MerchantRepository
}

func NewKafkaConsumer(merchantRepository model.MerchantRepository) *KafkaConsumer {
	kConsumer := &KafkaConsumer{merchantRepository: merchantRepository}
	kConsumer.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{os.Getenv("KAFKA_BROKER")},
		GroupID:   "register_merchant",
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

func (kc *KafkaConsumer) RegisterMerchant() {
	for {
		b := kc.Consume()
		log.Printf("Mensagem consumida %v\n", string(b))

		if len(b) > 0 {
			merchant := &model.Merchant{}
			err := json.Unmarshal(b, merchant)
			if err == nil {
				log.Printf("Salvando merchant em banco %v\n", merchant)
				err := kc.merchantRepository.Save(merchant)
				if err != nil {
					log.Printf("Error saving data %v\n", err)
				}

			} else {
				log.Printf("Error unmarshaling data = %v\n", err)
			}
		}
	}
}
