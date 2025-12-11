package formatter

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ErrorFormatter(errs validator.ValidationErrors) map[string][]string {
	errorMap := make(map[string][]string)

	for _, err := range errs {
		var message string
		tag := err.Tag()
		field := err.Field()
		param := err.Param()

		switch tag {

		case "required":
			message = fmt.Sprintf("%s is required", field)
		case "unique":
			message = fmt.Sprintf("%s is already taken", field)
		case "email":
			message = fmt.Sprintf("%s must be email format", field)
		case "gender":
			message = fmt.Sprintf("%s must be male or female", field)
		case "min":
			message = fmt.Sprintf("%s must be at least %s characters", field, param)
		case "max":
			message = fmt.Sprintf("%s must be no more than %s characters", field, param)
		default:
			message = fmt.Sprintf("%s is invalid", field)
		}

		errorMap[err.Field()] = append(errorMap[err.Field()], message)
	}

	return errorMap
}
