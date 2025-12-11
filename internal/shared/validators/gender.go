package validators

import (
	"github.com/go-playground/validator/v10"
)

func Gender(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "f" || value == "m" {
		return true
	}
	return false
}
