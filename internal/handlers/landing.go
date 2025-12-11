package handlers

import (
	"net/http"

	"github.com/dliluashvili/cowatchit/internal/templates"
)

func HandleLanding(w http.ResponseWriter, r *http.Request) {
	component := templates.LandingPage()

	component.Render(r.Context(), w)
}
