package provider

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/danilomarques1/godemo/gw/api/dto"
	"github.com/danilomarques1/godemo/gw/api/model"
	"github.com/danilomarques1/godemo/gw/api/response"
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
		apiErr := &dto.ApiErrorDto{}
		if err := json.NewDecoder(resp.Body).Decode(apiErr); err != nil {
			return nil, err
		}
		return nil, response.NewApiError(apiErr.Message, resp.StatusCode)
	}

	cob := &model.Cob{}
	if err := json.NewDecoder(resp.Body).Decode(cob); err != nil {
		return nil, err
	}
	return cob, nil
}

func (ip *ItauProvider) FindCob(token, txid string) (*model.Cob, error) {
	request, err := http.NewRequest(http.MethodGet, ip.providerUrl+"/"+txid, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", "Bearer "+token)
	resp, err := ip.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		apiErr := &dto.ApiErrorDto{}
		if err := json.NewDecoder(resp.Body).Decode(apiErr); err != nil {
			return nil, err
		}
		return nil, response.NewApiError(apiErr.Message, resp.StatusCode)
	}

	cob := &model.Cob{}
	if err := json.NewDecoder(resp.Body).Decode(cob); err != nil {
		return nil, err
	}
	return cob, nil
}

// TODO implement cancel cob
