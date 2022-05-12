package service

import (
	"os"

	"github.com/golang-jwt/jwt"
)

type Claim struct {
	ClientId     string
	ClientSecret string
	jwt.StandardClaims
}

const TOKEN_EXPIRATION_TIME = 300

func GetToken(clientId, clientSecret string) (string, error) {
	claims := &Claim{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: TOKEN_EXPIRATION_TIME,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}
