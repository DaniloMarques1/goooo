package response

import (
	"encoding/json"
	"net/http"

	"github.com/danilomarques1/godemo/gw/api/dto"
	"github.com/go-playground/validator/v10"
)

func RespondERR(w http.ResponseWriter, err error) {
	switch err.(type) {
	case *ApiError:
		apiError := err.(*ApiError)
		errDto := dto.ApiErrorDto{Message: apiError.Msg}
		RespondJSON(w, errDto, apiError.Status)
	case validator.ValidationErrors:
		errors := getValidationErrors(err)
		errDto := dto.ApiErrorDto{Message: "Validation error", Errors: errors}
		RespondJSON(w, errDto, http.StatusBadRequest)
	default:
		errDto := dto.ApiErrorDto{Message: err.Error()}
		RespondJSON(w, errDto, http.StatusInternalServerError)
	}
}

func RespondJSON(w http.ResponseWriter, body interface{}, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func getValidationErrors(errors error) []dto.Error {
	validationErrors := make([]dto.Error, 0)
	for _, err := range errors.(validator.ValidationErrors) {
		e := dto.Error{Field: err.Field(), RejectedValue: err.Value()}
		validationErrors = append(validationErrors, e)
	}
	return validationErrors
}
