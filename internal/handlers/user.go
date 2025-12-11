package handlers

import (
	"fmt"
	"net/http"

	"github.com/dliluashvili/cowatchit/internal/helpers"
	"github.com/dliluashvili/cowatchit/internal/models"
	"github.com/dliluashvili/cowatchit/internal/services"
	"github.com/dliluashvili/cowatchit/internal/shared/constants"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(us *services.UserService) *UserHandler {
	return &UserHandler{
		userService: us,
	}
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(constants.SessionContextKey).(*models.Session)

	user, err := h.userService.Me(session.User.ID)

	if err != nil {
		fmt.Println("Error /me", err)
		helpers.SendJson(w, &helpers.Response{
			Data:    nil,
			Message: "Unable to fetch authed user",
			Status:  500,
		})
		return
	}

	helpers.SendJson(w, &helpers.Response{
		Data: map[string]any{
			"user": user,
		},
		Message: "all goood",
		Status:  200,
	})
}
