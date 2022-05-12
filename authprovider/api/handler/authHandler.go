package handler

import (
	"log"
	"net/http"

	"github.com/danilomarques1/godemo/authprovider/api/dto"
	"github.com/danilomarques1/godemo/authprovider/api/service"
	"github.com/danilomarques1/godemo/authprovider/api/util"
	chi "github.com/go-chi/chi/v5"
)

type AuthHandler interface {
	GetToken(w http.ResponseWriter, r *http.Request)
}

type AuthHandlerImpl struct{}

func NewAuthHandlerImpl() *AuthHandlerImpl {
	return &AuthHandlerImpl{}
}

func (ah *AuthHandlerImpl) ConfigureRoutes(router *chi.Mux) {
	router.Post("/api/oauth/token", ah.GetToken)
}

func (ah *AuthHandlerImpl) GetToken(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	clientId := r.Form.Get("client_id")
	log.Printf("Client id received = %v\n", clientId)
	if len(clientId) == 0 {
		util.RespondERR(w, "Invalid client id", http.StatusBadRequest)
		return
	}

	clientSecret := r.Form.Get("client_secret")
	log.Printf("Client secret received = %v\n", clientSecret)
	if len(clientSecret) == 0 {
		util.RespondERR(w, "Invalid client secret", http.StatusBadRequest)
		return
	}

	tokenStr, err := service.GetToken(clientId, clientSecret)
	if err != nil {
		log.Printf("Error getting token = %v\n", err)
		util.RespondERR(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	token := dto.Token{AccessToken: tokenStr, ExpiresIn: service.TOKEN_EXPIRATION_TIME}
	util.RespondJSON(w, token, http.StatusCreated)
	return
}
