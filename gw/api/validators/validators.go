package validators

import (
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"reflect"
)

// validates the key field based on the key type
// where a NATIONALID should have 11 characters
// where a MOBILEPHOBE should have 11 characters
// where a MERCHANTNATIONALID should have 16 characters
// where a RANDOMKEY should be a valid uuid
func ValidatePixKey(fl validator.FieldLevel) bool {
	log.Printf("Adding Pix key validation\n")
	key := fl.Parent().FieldByName("KeyType").String()
	pixKey := fl.Field().String()
	switch key {
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
		return false
	}
	return true
}

// return the field name inside the json tag
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
