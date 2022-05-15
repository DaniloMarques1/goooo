package service

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/danilomarques1/godemo/gw/api/dto"
	"github.com/danilomarques1/godemo/gw/api/response"
	"github.com/go-redis/redis/v8"
)

type TokenService interface {
	GetToken() (*dto.Token, error)
}

type TokenServiceImpl struct {
	client    *http.Client
	tokenUrl  string
	redisConn *redis.Client
}

func NewTokenServiceImpl(redisConn *redis.Client) *TokenServiceImpl {
	client := &http.Client{}
	tokenUrl := os.Getenv("AUTH_BASE_URL")
	return &TokenServiceImpl{client: client, tokenUrl: tokenUrl, redisConn: redisConn}
}

func (ts *TokenServiceImpl) GetToken() (*dto.Token, error) {
	token, err := ts.getFromCache()
	if err == nil {
		log.Printf("Getting token from cache...\n")
		return token, nil
	}
	log.Printf("Token not found on cache, requesting auth provider")

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

		ts.saveCache(token)
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

func (ts *TokenServiceImpl) getFromCache() (*dto.Token, error) {
	str, err := ts.redisConn.Get(context.Background(), "token").Result()
	if err != nil {
		return nil, err
	}
	var token dto.Token
	if err := json.Unmarshal([]byte(str), &token); err != nil {
		return nil, err
	}
	return &token, nil
}

func (ts *TokenServiceImpl) saveCache(token dto.Token) error {
	bytes, err := json.Marshal(token)
	if err != nil {
		return err
	}

	err = ts.redisConn.SetEX(context.Background(), "token", string(bytes), time.Second*time.Duration(token.ExpiresIn)).Err()
	if err != nil {
		log.Printf("Error = %v\n", err)
		return err
	}
	return nil
}
