package main

import (
	"log"
	"os"

	"github.com/danilomarques1/godemo/authprovider/api"
	"github.com/danilomarques1/godemo/authprovider/api/handler"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(os.Getenv("PORT"))
	authHandler := handler.NewAuthHandlerImpl()
	authHandler.ConfigureRoutes(server.Router)
	server.Start()
}
