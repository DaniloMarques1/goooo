package util

import (
	"encoding/json"
	"net/http"

	"github.com/danilomarques1/godemo/authprovider/api/dto"
)

func RespondERR(w http.ResponseWriter, msg string, status int) {
	err := dto.Error{Message: msg}
	RespondJSON(w, err, status)
}

func RespondJSON(w http.ResponseWriter, body interface{}, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
