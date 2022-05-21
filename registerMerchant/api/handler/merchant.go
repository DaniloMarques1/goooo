package handler

import (
	"net/http"

	"github.com/danilomarques1/godemo/registerMerchants/api/service"
)

type MerchantHandler struct {
	merchantService service.MerchantService
}

func NewMerchantService(merchantService service.MerchantService) *MerchantHandler {
	return &MerchantHandler{merchantService: merchantService}
}

func (mh *MerchantHandler) FindAll(w http.ResponseWriter, r *http.Request) {
}

func (mh *MerchantHandler) FindById(w http.ResponseWriter, r *http.Request) {
}
