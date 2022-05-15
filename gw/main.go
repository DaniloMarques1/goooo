package main

import (
	"context"
	"log"
	"os"

	"github.com/danilomarques1/godemo/gw/api"
	"github.com/danilomarques1/godemo/gw/api/cache"
	"github.com/danilomarques1/godemo/gw/api/dto"
	"github.com/danilomarques1/godemo/gw/api/handler"
	"github.com/danilomarques1/godemo/gw/api/producer"
	"github.com/danilomarques1/godemo/gw/api/provider"
	"github.com/danilomarques1/godemo/gw/api/repository"
	"github.com/danilomarques1/godemo/gw/api/service"
	"github.com/danilomarques1/godemo/gw/api/validators"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(os.Getenv("PORT"))
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}

	cobRepository := repository.NewCobMongoRepository(client, "cob")

	redisConn := redis.NewClient(&redis.Options{
		Addr: "0.0.0.0:6379",
		DB:   0,
	})
	if err := redisConn.Ping(context.Background()).Err(); err != nil {
		log.Fatal(err)
	}

	redisCache := cache.NewRedisCache[*dto.Token](redisConn)
	tokenService := service.NewTokenServiceImpl(redisCache)

	validate := validator.New()
	validate.RegisterTagNameFunc(validators.GetJsonTagName)
	validate.RegisterValidation("pix-key", validators.ValidatePixKey)

	itauProvider := provider.NewItauProvider()
	kafkaProducer, err := producer.NewKafkaProducer()
	if err != nil {
		log.Fatal(err)
	}

	cobHandler := handler.NewCobHandler(cobRepository, tokenService, validate, itauProvider, kafkaProducer)
	cobHandler.ConfigureRoutes(server.Router)

	server.Start()
}
