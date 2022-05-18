package main

import (
	"log"
	"os"

	"github.com/danilomarques1/godemo/registerMerchant/api"
	"github.com/danilomarques1/godemo/registerMerchant/api/consumer"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	kafkaConsumer := consumer.NewKafkaConsumer()

	// @@@
	go func () {
		for {
			b := kafkaConsumer.Consume()
			log.Printf("%v\n", string(b))
		}
	}()

	server := api.NewServer(os.Getenv("PORT"))
	server.Start()
}
