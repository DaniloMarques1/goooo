package dto

type Token struct {
	ExpiresIn   int64  `json:"expires_in"`
	AccessToken string `json:"access_token"`
}
