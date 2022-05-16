package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/danilomarques1/godemo/provider/api/dto"
	"github.com/danilomarques1/godemo/provider/api/model"
	"github.com/danilomarques1/godemo/provider/api/util"
	chi "github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CobHandler struct {
	cobRepository model.CobRepository
	validator     *validator.Validate
}

func NewCobHandler(cobRepository model.CobRepository, validator *validator.Validate) *CobHandler {
	return &CobHandler{cobRepository: cobRepository, validator: validator}
}

func (ch *CobHandler) ConfigureRoutes(router *chi.Mux) {
	router.Post("/cob", ch.CreateCob)
	router.Get("/cob/{txid}", ch.FindCob)
}

func (ch *CobHandler) CreateCob(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request arrived\n")
	var createCobDto dto.CreateCobDto
	if err := json.NewDecoder(r.Body).Decode(&createCobDto); err != nil {
		log.Printf("ERR parsing json = %v\n", err)
		apiErr := util.NewApiError("Invalid body", http.StatusBadRequest)
		util.RespondERR(w, apiErr)
		return
	}

	if err := ch.validator.Struct(createCobDto); err != nil {
		log.Printf("Error validating struct = %v\n", err)
		util.RespondERR(w, err)
		return
	}

	log.Printf("CreateCob = %v\n", createCobDto)

	cal := model.Calendar{
		CreatedAt: time.Now(),
		ExpiresIn: createCobDto.Cal.ExpiresIn,
	}
	addInfosDto := createCobDto.AdditionalInfos
	addInfos := make([]model.AdditionalInfo, 0)
	for _, info := range addInfosDto {
		addInfo := model.AdditionalInfo{Key: info.Key, Value: info.Value}
		addInfos = append(addInfos, addInfo)
	}

	cob := &model.Cob{
		TxId:            uuid.NewString(),
		Value:           createCobDto.Value,
		KeyType:         createCobDto.KeyType,
		Key:             createCobDto.Key,
		Cal:             cal,
		AdditionalInfos: addInfos,
	}
	log.Printf("Calling repository = %v\n", cob)

	if err := ch.cobRepository.Save(cob); err != nil {
		log.Printf("ERR saving cob = %v\n", err)
		util.RespondERR(w, err)
		return
	}
	log.Printf("Cob saved\n")

	util.RespondJSON(w, cob, http.StatusCreated)
	return
}

func (ch *CobHandler) FindCob(w http.ResponseWriter, r *http.Request) {
	log.Printf("Consult cob\n")
	txid := chi.URLParam(r, "txid")
	if _, err := uuid.Parse(txid); err != nil {
		apiErr := util.NewApiError("Invalid txid", http.StatusBadRequest)
		util.RespondERR(w, apiErr)
		return
	}
	cob, err := ch.cobRepository.FindById(txid)
	if err != nil {
		util.RespondERR(w, err)
		return
	}
	util.RespondJSON(w, cob, http.StatusOK)
	return
}
