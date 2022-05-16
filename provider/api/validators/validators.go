package validators

import (
	"log"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func GetJsonTagName(fld reflect.StructField) string {
	splitted := strings.SplitN(fld.Tag.Get("json"), ",", 2)
	if len(splitted) == 0 {
		return ""
	}

	name := splitted[0]
	if name == "-" {
		return ""
	}

	return name
}

func ValidatePixKey(fl validator.FieldLevel) bool {
	log.Printf("Adding Pix key validation\n")
	parent := fl.Parent()
	pixKey := fl.Field().String()
	switch parent.FieldByName("KeyType").String() {
	case "NATIONALID":
	case "MOBILEPHONE":
		return len(pixKey) == 11
	case "MERCHANTNATIONALID":
		return len(pixKey) == 16
	case "RANDOMKEY":
		if _, err := uuid.Parse(pixKey); err != nil {
			return false
		}
		return true

	default:
		return true
	}
	return true
}
