package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func Username(fl validator.FieldLevel) bool {
	username := fl.Field().String()

	length := len(username)

	if length < 4 || length > 15 {
		return false
	}

	if !usernameRegex.MatchString(username) {
		return false
	}

	return true
}
