package validators

import (
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

var dobRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

func Dob(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	if len(value) == 0 {
		return false
	}

	if !dobRegex.MatchString(value) {
		return false
	}

	date, err := time.Parse("2006-01-02", value)

	if err != nil {
		return false
	}

	currentYear := time.Now().Year()
	minYear := currentYear - 100
	maxYear := 2007

	year := date.Year()

	if year < minYear || year > maxYear {
		return false
	}

	return true
}
