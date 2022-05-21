package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/danilomarques1/godemo/gw/api/dto"
	"github.com/danilomarques1/godemo/gw/api/model"
	"github.com/danilomarques1/godemo/gw/api/producer"
	"github.com/danilomarques1/godemo/gw/api/validators"
	"github.com/go-playground/validator/v10"
)

type CobRepositoryMock struct{}

func (crm *CobRepositoryMock) Save(cob *model.Cob) error {
	return nil
}

func (crm *CobRepositoryMock) FindById(txid string) (*model.Cob, error) {
	return nil, nil
}

type TokenServiceMock struct{}

func (tsm *TokenServiceMock) GetToken() (*dto.Token, error) {
	return &dto.Token{
		AccessToken: "access_token",
		ExpiresIn:   time.Now().Unix() + 300,
	}, nil
}

type ProviderMock struct{}

func (pm *ProviderMock) CreateCob(token string, cobDto dto.CreateCobDto) (*model.Cob, error) {
	return nil, nil
}

func (pm *ProviderMock) FindCob(token, txid string) (*model.Cob, error) {
	return nil, nil
}

type ProducerMock struct{}

func (pm *ProducerMock) Produce(merchant producer.Merchant) error {
	return nil
}

func (pm *ProducerMock) Close() error {
	return nil
}

func TestCreateCob(t *testing.T) {
	body := `
		{	
			"calendar": {
				"expires_in": 120
			},
			"value": 8.30,
			"key": "11954769490",
			"key_type": "NATIONALID",
			"additional_info": [
				{
					"key": "sub_acquirer_id",
					"value": "31"
				},
				{
					"key": "sub_acquirer_name",
					"value": "Phoebus Team Dev"
				},
				{
					"key": "merchant_id",
					"value": "000200"
				},
				{
					"key": "merchant_name",
					"value": "Jacuma"
				},
				{
					"key": "merchant_address",
					"value": "Recife"
				},
				{
					"key": "terminal_d",
					"value": "3020018"
				},
				{
					"key": "app_version",
					"value": "v1.0.0"
				}
			]
		}
	`
	request := httptest.NewRequest(http.MethodPost, "/cob", strings.NewReader(body))
	rr := httptest.NewRecorder()
	validate := validator.New()
	validate.RegisterTagNameFunc(validators.GetJsonTagName)
	validate.RegisterValidation("pix-key", validators.ValidatePixKey)

	handler := NewCobHandler(&CobRepositoryMock{}, &TokenServiceMock{}, validate, &ProviderMock{}, &ProducerMock{})

	handler.CreateCob(rr, request)

	if rr.Code != http.StatusCreated {
		t.Fatalf("Wrong status code returned. Expected %v received %v\n",
			http.StatusCreated, rr.Code)
	}
}
