package api

import (
	"fmt"
	"log"
	"net/http"

	chi "github.com/go-chi/chi/v5"
)

type Server struct {
	Router *chi.Mux
	port   string
}

func NewServer(port string) *Server {
	s := &Server{}
	s.Router = chi.NewRouter()
	s.port = port
	return s
}

func (s *Server) Start() {
	log.Printf("Starting server on port %v\n", s.port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", s.port), s.Router))
}
