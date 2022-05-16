package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/danilomarques1/godemo/provider/api/service"
	"github.com/danilomarques1/godemo/provider/api/util"
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
	log.Printf("Server starting on port %v\n", s.port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", s.port), s.Router))
}

func (s *Server) Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getTokenFromHeader(r.Header)
		if err != nil {
			util.RespondERR(w, err)
			return
		}
		err = service.ValidateToken(token)
		if err != nil {
			util.RespondERR(w, err)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getTokenFromHeader(header http.Header) (string, error) {
	bearer := header.Get("Authorization")
	splitted := strings.Split(bearer, " ")
	if len(splitted) < 2 {
		return "", util.NewApiError("Token not provided", http.StatusUnauthorized)
	}
	return splitted[1], nil
}
