package handlers

import (
	"fmt"
	"net/http"

	"github.com/dliluashvili/cowatchit/internal/dtos"
	"github.com/dliluashvili/cowatchit/internal/helpers"
	"github.com/dliluashvili/cowatchit/internal/services"
	"github.com/dliluashvili/cowatchit/internal/shared/constants"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(s *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: s,
	}
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	validated := r.Context().Value(constants.ValidatedContextKey).(*dtos.SignInDto)

	sessionModel, code, err := h.authService.SignIn(r.Context(), validated)

	if err != nil {
		helpers.SendJson(w, &helpers.Response{
			Data: map[string]bool{
				"success": false,
			},
			Status:  code,
			Message: "Invalid Credentials",
		})

		return
	}

	sessionID := sessionModel.SessionID

	http.SetCookie(w, &http.Cookie{
		Name:     "sessionId",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   int(constants.SessionDuration.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	helpers.SendJson(w, &helpers.Response{
		Data: map[string]bool{
			"success": true,
		},
		Message: "User has been created",
		Status:  code,
	})
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	validated := r.Context().Value(constants.ValidatedContextKey).(*dtos.SignUpDto)

	sessionModel, code, err := h.authService.SignUp(r.Context(), validated)

	if err != nil {
		fmt.Println("error while signing up", err)
		helpers.SendJson(w, &helpers.Response{
			Data: map[string]bool{
				"success": false,
			},
			Status:  code,
			Message: "Strange error, try again",
		})

		return
	}

	sessionID := sessionModel.SessionID

	http.SetCookie(w, &http.Cookie{
		Name:     "sessionId",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   int(constants.SessionDuration.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	helpers.SendJson(w, &helpers.Response{
		Data: map[string]bool{
			"success": true,
		},
		Message: "User has been created",
		Status:  code,
	})

}
