package service

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/danilomarques1/godemo/gw/api/cache"
	"github.com/danilomarques1/godemo/gw/api/dto"
	"github.com/danilomarques1/godemo/gw/api/response"
)

type TokenService interface {
	GetToken() (*dto.Token, error)
}

type TokenServiceImpl struct {
	client        *http.Client
	tokenUrl      string
	cacheProvider cache.Cache[*dto.Token]
}

func NewTokenServiceImpl(cacheProvider cache.Cache[*dto.Token]) *TokenServiceImpl {
	client := &http.Client{}
	tokenUrl := os.Getenv("AUTH_BASE_URL")
	return &TokenServiceImpl{client: client, tokenUrl: tokenUrl, cacheProvider: cacheProvider}
}

func (ts *TokenServiceImpl) GetToken() (*dto.Token, error) {
	cachedToken := &dto.Token{}
	err := ts.cacheProvider.GetFromCache("token", cachedToken)
	if err == nil {
		log.Printf("Getting token from cache...\n")
		return cachedToken, nil
	}
	log.Printf("Token not found on cache, requesting auth provider %v\n", err)

	token, err := ts.requestAuthProvider()
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (ts *TokenServiceImpl) requestAuthProvider() (*dto.Token, error) {
	body := url.Values{}
	body.Add("client_id", os.Getenv("CLIENT_ID"))
	body.Add("client_secret", os.Getenv("CLIENT_SECRET"))

	request, err := http.NewRequest(http.MethodPost, ts.tokenUrl, strings.NewReader(body.Encode()))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := ts.client.Do(request)
	if err != nil {
		log.Printf("Error doing request = %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		var token dto.Token
		if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
			log.Printf("Error decoding json = %v\n", err)
			return nil, err
		}

		ts.cacheProvider.SaveToCache("token", &token, time.Second*time.Duration(token.ExpiresIn))
		return &token, nil
	} else if resp.StatusCode == http.StatusBadRequest {
		var authError dto.AuthResponseError
		if err := json.NewDecoder(resp.Body).Decode(&authError); err != nil {
			log.Printf("Error decoding error json = %v\n", err)
			return nil, err
		}
		return nil, response.NewApiError(authError.Message, http.StatusBadRequest)
	} else {
		return nil, response.NewApiError("Internal server error", http.StatusInternalServerError)
	}
}
