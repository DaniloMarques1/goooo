package main

import (
	"context"
	"log"
	"os"

	"github.com/danilomarques1/godemo/provider/api"
	"github.com/danilomarques1/godemo/provider/api/handler"
	"github.com/danilomarques1/godemo/provider/api/repository"
	"github.com/danilomarques1/godemo/provider/api/validators"
	"github.com/go-playground/validator/v10"
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
	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal("Impossible to ping %v\n", err)
	}

	cobRepository := repository.NewCobMongoRepository(client, "cob")
	validate := validator.New()
	validate.RegisterTagNameFunc(validators.GetJsonTagName)
	validate.RegisterValidation("pix-key", validators.ValidatePixKey)
	cobHandler := handler.NewCobHandler(cobRepository, validate)
	cobHandler.ConfigureRoutes(server.Router)

	server.Start()
}
