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
	s.port = port
	s.Router = chi.NewRouter()
	s.Router.Use(applicationJson)
	return s
}

func (s *Server) Start() {
	log.Printf("Server starting at port %v\n", s.port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", s.port), s.Router))
}

func applicationJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
