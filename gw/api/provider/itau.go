package provider

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/danilomarques1/godemo/gw/api/dto"
	"github.com/danilomarques1/godemo/gw/api/model"
	"github.com/danilomarques1/godemo/gw/api/util"
)

type ItauProvider struct {
	client      *http.Client
	providerUrl string
}

func NewItauProvider() *ItauProvider {
	providerUrl := os.Getenv("PROVIDER_URL")
	client := &http.Client{}
	return &ItauProvider{client: client, providerUrl: providerUrl}
}

func (ip *ItauProvider) CreateCob(token string, cobDto dto.CreateCobDto) (*model.Cob, error) {
	b, err := json.Marshal(&cobDto)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, ip.providerUrl, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", "Bearer "+token)
	resp, err := ip.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var apiError dto.ApiErrorDto
		if err := json.NewDecoder(resp.Body).Decode(&apiError); err != nil {
			return nil, err
		}
		return nil, util.NewApiError(apiError.Message, resp.StatusCode)
	}

	var cobResp model.Cob
	if err := json.NewDecoder(resp.Body).Decode(&cobResp); err != nil {
		return nil, err
	}
	return &cobResp, nil
}

func (ip *ItauProvider) FindCob(token, txid string) (*model.Cob, error) {
	return nil, nil
}
