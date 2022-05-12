package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/danilomarques1/godemo/authprovider/api/handler"
	chi "github.com/go-chi/chi/v5"
)

type Server struct {
	Router      *chi.Mux
	port        int64
	authHandler handler.AuthHandler
}

func NewServer(port int64) *Server {
	s := &Server{}
	s.port = port
	s.Router = chi.NewRouter()
	return s
}

func (s *Server) Start() {
	log.Printf("Server running on port %v\n", s.port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", s.port), s.Router))
}
