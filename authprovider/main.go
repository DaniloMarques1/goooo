package main

import (
	"log"

	"github.com/danilomarques1/godemo/authprovider/api"
	"github.com/danilomarques1/godemo/authprovider/api/handler"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(8080)
	authHandler := handler.NewAuthHandlerImpl()
	authHandler.ConfigureRoutes(server.Router)
	server.Start()
}
