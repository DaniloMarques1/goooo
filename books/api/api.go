package api

import (
	"log"
	"net/http"

	chi "github.com/go-chi/chi/v5"
)

type server struct {
	Router chi.Router
	port   string
}

func NewServer(port string) *server {
	s := &server{port: port}
	s.Router = chi.NewRouter()
	return s
}

func (s *server) Start() {
	log.Fatal(http.ListenAndServe(":"+s.port, s.Router))
}
