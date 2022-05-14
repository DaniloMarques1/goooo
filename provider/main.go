package main

import (
	"context"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/danilomarques1/godemo/provider/api"
	"github.com/danilomarques1/godemo/provider/api/handler"
	"github.com/danilomarques1/godemo/provider/api/repository"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(8081)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal("Impossible to ping %v\n", err)
	}

	cobRepository := repository.NewCobMongoRepository(client, "cob")
	validator := validator.New()
	validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
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
	cobHandler := handler.NewCobHandler(cobRepository, validator)
	cobHandler.ConfigureRoutes(server.Router)

	server.Start()
}
