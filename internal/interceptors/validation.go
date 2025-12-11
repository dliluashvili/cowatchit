package interceptors

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dliluashvili/cowatchit/internal/helpers"
	"github.com/dliluashvili/cowatchit/internal/shared/constants"
	"github.com/dliluashvili/cowatchit/internal/shared/formatter"
	"github.com/go-playground/validator/v10"
)

func ValidateBody[T any](validate *validator.Validate) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var dto T

			if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			if err := validate.Struct(dto); err != nil {
				if ve, ok := err.(validator.ValidationErrors); ok {
					formatted := formatter.ErrorFormatter(ve)

					helpers.SendJson(w, &helpers.Response{
						Data:    formatted,
						Message: "Validation failed",
						Status:  http.StatusUnprocessableEntity,
					})

					return
				}

				http.Error(w, "Validation failed", http.StatusBadRequest)
				return
			}

			ctx := context.WithValue(r.Context(), constants.ValidatedContextKey, &dto)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
