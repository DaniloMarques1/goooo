package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/danilomarques1/godemo/provider/api/dto"
	"github.com/danilomarques1/godemo/provider/api/model"
	chi "github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CobHandler struct {
	cobRepository model.CobRepository
}

func NewCobHandler(cobRepository model.CobRepository) *CobHandler {
	return &CobHandler{cobRepository: cobRepository}
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
		return // TODO return some error
	}
	log.Printf("CreateCob = %v\n", createCobDto)

	cal := model.Calendar{
		CreatedAt: time.Now(),
		ExpiresIn: createCobDto.Cal.ExpiresIn,
	}
	addInfosDto := createCobDto.AdditionalInfos
	addInfos := make([]model.AdditionalInfo, len(addInfosDto))
	for _, info := range addInfosDto {
		addInfo := model.AdditionalInfo{Key: info.Key, Value: info.Value}
		addInfos = append(addInfos, addInfo)
	}

	cob := &model.Cob{
		TxId:            uuid.NewString(),
		Value:           createCobDto.Value,
		Cal:             cal,
		AdditionalInfos: addInfos,
	}
	log.Printf("Calling repository = %v\n", cob)

	if err := ch.cobRepository.Save(cob); err != nil {
		log.Printf("ERR saving cob = %v\n", err)
		return // TODO return some error
	}
	log.Printf("Cob saved\n")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cob)
	return
}

func (ch *CobHandler) FindCob(w http.ResponseWriter, r *http.Request) {
	return
}
