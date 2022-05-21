package handler

import (
	"encoding/json"
	"net/http"

	"github.com/danilomarques1/godemo/registerMerchant/api/model"
	chi "github.com/go-chi/chi/v5"
)

type MerchantHandler struct {
	merchantRepository model.MerchantRepository
}

func NewMerchantHandler(merchantRepository model.MerchantRepository) *MerchantHandler {
	return &MerchantHandler{merchantRepository: merchantRepository}
}

type ErrorDto struct {
	Message string `json:"message"`
}

type MerchantsDto struct {
	Merchants []model.Merchant `json:"merchants"`
}

type MerchantDto struct {
	Merchant *model.Merchant `json:"merchant"`
}

func (mh *MerchantHandler) ConfigureRoutes(router chi.Router) {
	router.Get("/merchants", mh.FindAll)
	router.Get("/merchants/{merchant_id}", mh.FindById)
}

func (mh *MerchantHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	merchants, err := mh.merchantRepository.FindAll()
	if err != nil {
		resp := ErrorDto{Message: err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	json.NewEncoder(w).Encode(MerchantsDto{Merchants: merchants})
	return
}

func (mh *MerchantHandler) FindById(w http.ResponseWriter, r *http.Request) {
	merchantId := chi.URLParam(r, "merchant_id")
	merchant, err := mh.merchantRepository.FindById(merchantId)
	if err != nil {
		resp := ErrorDto{Message: err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	json.NewEncoder(w).Encode(MerchantDto{Merchant: merchant})
	return
}
