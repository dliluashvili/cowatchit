package validators

import (
	"time"

	"github.com/go-playground/validator/v10"
)

func Date(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	_, err := time.Parse("2006-01-02", value)

	return err == nil
}
