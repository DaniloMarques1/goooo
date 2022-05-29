package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/danilomarques1/godemo/gw/api/dto"
	"github.com/danilomarques1/godemo/gw/api/model"
	"github.com/danilomarques1/godemo/gw/api/producer"
	"github.com/danilomarques1/godemo/gw/api/provider"
	"github.com/danilomarques1/godemo/gw/api/response"
	"github.com/danilomarques1/godemo/gw/api/service"
	chi "github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CobHandler struct {
	cobRepository model.CobRepository
	tokenService  service.TokenService
	prov          provider.Provider
	validate      *validator.Validate
	prod          producer.Producer
}

func NewCobHandler(cobRepository model.CobRepository, tokenService service.TokenService,
	validate *validator.Validate, prov provider.Provider, prod producer.Producer) *CobHandler {
	return &CobHandler{
		cobRepository: cobRepository,
		tokenService:  tokenService,
		validate:      validate,
		prov:          prov,
		prod:          prod,
	}
}

func (ch *CobHandler) ConfigureRoutes(router chi.Router) {
	router.Post("/cob", ch.CreateCob)
	router.Get("/cob/{txid}", ch.FindCob)
}

func (ch *CobHandler) CreateCob(w http.ResponseWriter, r *http.Request) {
	var cobCreateDto dto.CreateCobDto
	if err := json.NewDecoder(r.Body).Decode(&cobCreateDto); err != nil {
		apiErr := response.NewApiError("Invalid body", http.StatusBadRequest)
		response.RespondERR(w, apiErr)
		return
	}
	if err := ch.validate.Struct(cobCreateDto); err != nil {
		response.RespondERR(w, err)
		return
	}
	token, err := ch.tokenService.GetToken()
	if err != nil {
		response.RespondERR(w, err)
		return
	}

	resp, err := ch.prov.CreateCob(token.AccessToken, cobCreateDto)
	if err != nil {
		response.RespondERR(w, err)
		return
	}
	go ch.produceMessage(resp) // no need for waiting

	if err := ch.cobRepository.Save(resp); err != nil {
		response.RespondERR(w, err)
		return
	}

	response.RespondJSON(w, resp, http.StatusCreated)
	return
}

func (ch *CobHandler) FindCob(w http.ResponseWriter, r *http.Request) {
	txid := chi.URLParam(r, "txid")
	if _, err := uuid.Parse(txid); err != nil {
		apiErr := response.NewApiError("Invalid txid", http.StatusBadRequest)
		response.RespondERR(w, apiErr)
		return
	}
	cob, err := ch.cobRepository.FindById(txid)
	if err != nil {
		response.RespondERR(w, err)
		return
	}

	response.RespondJSON(w, cob, http.StatusOK)
	return
}

func (ch *CobHandler) produceMessage(cob *model.Cob) {
	addInfos := cob.AdditionalInfos
	if len(addInfos) == 0 {
		return
	}

	var merchant producer.Merchant
	for _, info := range addInfos {
		switch info.Key {
		case "sub_acquirer_id":
			merchant.SubAcquirerId = info.Value
		case "sub_acquirer_name":
			merchant.SubAcquirerName = info.Value
		case "merchant_id":
			merchant.MerchantId = info.Value
		case "merchant_name":
			merchant.MerchantName = info.Value
		case "merchant_address":
			merchant.MerchantAddress = info.Value
		}
	}

	if err := ch.prod.Produce(merchant); err != nil {
		log.Printf("Error producing message = %v\n", err)
	}
}
