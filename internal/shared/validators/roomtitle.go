package validators

import (
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var roomTitleRegex = regexp.MustCompile(`^[\p{L}\p{N}\p{P}\s]+$`)
var consecutiveSpacesRegex = regexp.MustCompile(`\s{2,}`)

func RoomTitle(fl validator.FieldLevel) bool {
	title := fl.Field().String()

	// Check length constraints
	length := len(title)
	if length < 1 || length > 200 {
		return false
	}

	// Check for multiple consecutive spaces
	if consecutiveSpacesRegex.MatchString(title) {
		return false
	}

	// Check if title contains at least one meaningful character
	hasContent := false
	for _, r := range title {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			hasContent = true
			break
		}
	}
	if !hasContent {
		return false
	}

	// Check against Unicode-aware regex
	if !roomTitleRegex.MatchString(title) {
		return false
	}

	return true
}
