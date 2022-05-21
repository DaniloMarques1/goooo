package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/danilomarques1/godemo/registerMerchant/api"
	"github.com/danilomarques1/godemo/registerMerchant/api/consumer"
	"github.com/danilomarques1/godemo/registerMerchant/api/handler"
	"github.com/danilomarques1/godemo/registerMerchant/api/repository"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const schema = `
CREATE TABLE IF NOT EXISTS merchant(
  merchant_id varchar(30) primary key,
  merchant_name varchar(50) not null,
  merchant_address varchar(50) not null,
  sub_acquirer_id  varchar(30) not null,
  sub_acquirer_name varchar(50) not null
);
`

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	dataSource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"),
		os.Getenv("DB_PWD"), os.Getenv("DB_NAME"))
	db, err := sql.Open("postgres", dataSource)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(schema); err != nil {
		log.Fatal(err)
	}

	merchantRepository := repository.NewMerchantSqlRepository(db)
	kafkaConsumer := consumer.NewKafkaConsumer(merchantRepository)
	go kafkaConsumer.RegisterMerchant()

	server := api.NewServer(os.Getenv("PORT"))

	merchantHandler := handler.NewMerchantHandler(merchantRepository)
	merchantHandler.ConfigureRoutes(server.Router)

	server.Start()
}
