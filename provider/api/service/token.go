package service

import (
	"net/http"
	"os"

	"github.com/danilomarques1/godemo/provider/api/util"
	"github.com/golang-jwt/jwt"
)

func ValidateToken(tokenString string) error {
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})
	if err != nil {
		return util.NewApiError("Token is invalid", http.StatusUnauthorized)
	}
	return nil
}
