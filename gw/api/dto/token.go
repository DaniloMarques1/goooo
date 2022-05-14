package dto

type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type AuthResponseError struct {
	Message string `json:"message"`
}
