package handler

import (
	"encoding/json"
	"net/http"

	"github.com/danilomarques1/godemo/gw/api/dto"
	"github.com/danilomarques1/godemo/gw/api/model"
	"github.com/danilomarques1/godemo/gw/api/provider"
	"github.com/danilomarques1/godemo/gw/api/service"
	"github.com/danilomarques1/godemo/gw/api/util"
	chi "github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type CobHandler struct {
	cobRepository model.CobRepository
	tokenService  service.TokenService
	prov          provider.Provider
	validate      *validator.Validate
}

func NewCobHandler(cobRepository model.CobRepository, tokenService service.TokenService,
	validate *validator.Validate, prov provider.Provider) *CobHandler {
	return &CobHandler{
		cobRepository: cobRepository,
		tokenService:  tokenService,
		validate:      validate,
		prov:          prov,
	}
}

func (ch *CobHandler) ConfigureRoutes(router *chi.Mux) {
	router.Post("/cob", ch.CreateCob)
}

func (ch *CobHandler) CreateCob(w http.ResponseWriter, r *http.Request) {
	var cobCreateDto dto.CreateCobDto
	if err := json.NewDecoder(r.Body).Decode(&cobCreateDto); err != nil {
		util.RespondERR(w, err)
		return
	}
	if err := ch.validate.Struct(cobCreateDto); err != nil {
		util.RespondERR(w, err)
		return
	}
	token, err := ch.tokenService.GetToken()
	if err != nil {
		util.RespondERR(w, err)
		return
	}

	resp, err := ch.prov.CreateCob(token.AccessToken, cobCreateDto)
	if err != nil {
		util.RespondERR(w, err)
		return
	}

	if err := ch.cobRepository.Save(resp); err != nil {
		util.RespondERR(w, err)
		return
	}

	util.RespondJSON(w, resp, http.StatusCreated)
	return
}
