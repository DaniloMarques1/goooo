package main

import (
	"context"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/danilomarques1/godemo/gw/api"
	"github.com/danilomarques1/godemo/gw/api/handler"
	"github.com/danilomarques1/godemo/gw/api/provider"
	"github.com/danilomarques1/godemo/gw/api/repository"
	"github.com/danilomarques1/godemo/gw/api/service"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
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

	tokenService := service.NewTokenServiceImpl(redisConn)
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		splitted := strings.SplitN(fld.Tag.Get("json"), ",", 2)
		if len(splitted) == 0 {
			return ""
		}

		name := splitted[0]
		if name == "-" {
			return ""
		}

		return name

	})
	validate.RegisterValidation("pix-key", ValidatePixKey)
	itauProvider := provider.NewItauProvider()
	cobHandler := handler.NewCobHandler(cobRepository, tokenService, validate, itauProvider)
	cobHandler.ConfigureRoutes(server.Router)

	server.Start()
}

func ValidatePixKey(fl validator.FieldLevel) bool {
	log.Printf("Adding Pix key validation\n")
	parent := fl.Parent()
	pixKey := fl.Field().String()
	switch parent.FieldByName("KeyType").String() {
	case "NATIONALID":
	case "MOBILEPHONE":
		return len(pixKey) == 11
	case "MERCHANTNATIONALID":
		return len(pixKey) == 16
	case "RANDOMKEY":
		if _, err := uuid.Parse(pixKey); err != nil {
			return false
		}
		return true

	default:
		return true
	}
	return true
}
