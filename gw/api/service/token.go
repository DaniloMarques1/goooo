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
	cacheProvider cache.Cache
}

func NewTokenServiceImpl(cacheProvider cache.Cache) *TokenServiceImpl {
	client := &http.Client{}
	tokenUrl := os.Getenv("AUTH_BASE_URL")
	return &TokenServiceImpl{client: client, tokenUrl: tokenUrl, cacheProvider: cacheProvider}
}

func (ts *TokenServiceImpl) GetToken() (*dto.Token, error) {
	cachedToken, err := ts.getFromCache()
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

func (ts *TokenServiceImpl) getFromCache() (*dto.Token, error) {
	b, err := ts.cacheProvider.GetFromCache("token")
	if err != nil {
		return nil, err
	}
	token := &dto.Token{}
	if err := json.Unmarshal(b, token); err != nil {
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

	token, err := ts.parseResponse(resp)
	if err != nil {
		log.Printf("Error parsing response %v\n", err)
		return nil, err
	}

	ts.saveTokenToCache(token)
	return token, nil
}

func (ts *TokenServiceImpl) saveTokenToCache(token *dto.Token) error {
	expiration := time.Second * time.Duration(token.ExpiresIn)
	b, err := json.Marshal(token)
	if err != nil {
		log.Printf("Error parsing token %v\n", err)
		return err
	}

	err = ts.cacheProvider.SaveToCache("token", b, expiration)
	if err != nil {
		log.Printf("Error saving token to cache %v\n", err)
		return err
	}
	return nil
}

func (ts *TokenServiceImpl) parseResponse(resp *http.Response) (*dto.Token, error) {
	if resp.StatusCode == http.StatusCreated {
		token := &dto.Token{}
		if err := json.NewDecoder(resp.Body).Decode(token); err != nil {
			log.Printf("Error decoding json = %v\n", err)
			return nil, err
		}

		return token, nil
	} else {
		authError := &dto.AuthResponseError{}
		if err := json.NewDecoder(resp.Body).Decode(authError); err != nil {
			log.Printf("Error decoding error json = %v\n", err)
			return nil, err
		}

		return nil, response.NewApiError(authError.Message, resp.StatusCode)
	}
}
