package handler

import (
	"bytes"
	"encoding/json"
	//"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/danilomarques1/godemo/gw/api/dto"
	"github.com/danilomarques1/godemo/gw/api/model"
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

func (crm *CobRepositoryMock) Update(cob *model.Cob) error {
	return nil
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
	return &model.Cob{
		TxId:            "cee8eda2-1baa-4c75-b58b-dca145fae385",
		Value:           10.0,
		Status:          model.ACTIVE,
		KeyType:         "RANDOMKEY",
		Key:             "f4ff91ab-aa4c-4fd0-ba38-0ec624bc2509",
		Cal:             model.Calendar{ExpiresIn: 300, CreatedAt: time.Now()},
		AdditionalInfos: []model.AdditionalInfo{{Key: "merchant_id", Value: "12345"}},
	}, nil
}

func (pm *ProviderMock) FindCob(token, txid string) (*model.Cob, error) {
	return &model.Cob{
		TxId:            txid,
		Value:           10.0,
		Status:          model.ACTIVE,
		KeyType:         "RANDOMKEY",
		Key:             "f4ff91ab-aa4c-4fd0-ba38-0ec624bc2509",
		Cal:             model.Calendar{ExpiresIn: 300, CreatedAt: time.Now()},
		AdditionalInfos: []model.AdditionalInfo{{Key: "merchant_id", Value: "12345"}},
	}, nil
}

func (pm *ProviderMock) Cancel(token, txid string) error {
	return nil
}

type ProducerMock struct{}

func (pm *ProducerMock) Produce(b []byte) error {
	return nil
}

func (pm *ProducerMock) Close() error {
	return nil
}

func TestCreateCob(t *testing.T) {
	cases := []struct {
		label          string
		body           *dto.CreateCobDto
		statusExpected int
	}{
		{
			label: "Should return http status 201",
			body: &dto.CreateCobDto{
				Value:           10.0,
				KeyType:         "RANDOMKEY",
				Key:             "f4ff91ab-aa4c-4fd0-ba38-0ec624bc2509",
				Cal:             dto.CalendarDto{ExpiresIn: 300},
				AdditionalInfos: []dto.AdditionalInfoDto{{Key: "merchant_id", Value: "123456"}},
			},
			statusExpected: http.StatusCreated,
		},
		{
			label: "Should return http status 400 missing required field",
			body: &dto.CreateCobDto{
				KeyType:         "RANDOMKEY",
				Key:             "f4ff91ab-aa4c-4fd0-ba38-0ec624bc2509",
				Cal:             dto.CalendarDto{ExpiresIn: 300},
				AdditionalInfos: []dto.AdditionalInfoDto{{Key: "merchant_id", Value: "123456"}},
			},
			statusExpected: http.StatusBadRequest,
		},
		{
			label:          "Should return http status 400 missing body",
			body:           nil,
			statusExpected: http.StatusBadRequest,
		},
	}

	validate := validator.New()
	validate.RegisterTagNameFunc(validators.GetJsonTagName)
	validate.RegisterValidation("pix-key", validators.ValidatePixKey)

	handler := NewCobHandler(&CobRepositoryMock{}, &TokenServiceMock{}, validate, &ProviderMock{}, &ProducerMock{})

	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			b, err := json.Marshal(tc.body)
			if err != nil {
				t.Fatalf("Error = %v\n", err)
			}
			request := httptest.NewRequest(http.MethodPost, "/cob", bytes.NewReader(b))
			rr := httptest.NewRecorder()
			handler.CreateCob(rr, request)

			if rr.Code != tc.statusExpected {
				t.Fatalf("Expected status %v instead got %v\n", tc.statusExpected, rr.Code)
			}
		})
	}
}
