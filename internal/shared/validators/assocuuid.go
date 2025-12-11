package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func AssocUUID(fl validator.FieldLevel) bool {
	data, ok := fl.Field().Interface().(map[string]string)
	if !ok {
		return false
	}

	for key, val := range data {
		if _, err := uuid.Parse(key); err != nil {
			return false
		}
		if _, err := uuid.Parse(val); err != nil {
			return false
		}
	}

	return true
}
