package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/dliluashvili/cowatchit/internal/helpers"
	"github.com/dliluashvili/cowatchit/internal/services"

	"github.com/dliluashvili/cowatchit/internal/shared/constants"
)

func AuthSession(sessionService *services.SessionService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionID := extractSessionID(r)
			if sessionID == "" {
				handleUnauthorized(w, r)
				return
			}

			session, err := sessionService.GetUserBySession(r.Context(), sessionID)

			if err != nil || session == nil {
				handleUnauthorized(w, r)
				return
			}

			// Inject session into context
			ctx := context.WithValue(r.Context(), constants.SessionContextKey, session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractSessionID(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}

	if sessionID := r.Header.Get("x-session-id"); sessionID != "" {
		return sessionID
	}

	cookie, err := r.Cookie("sessionId")
	if err == nil {
		return cookie.Value
	}

	return ""
}

func handleUnauthorized(w http.ResponseWriter, r *http.Request) {
	isAJAX := r.Header.Get("X-Requested-With") == "XMLHttpRequest" ||
		r.Header.Get("HX-Request") == "true" ||
		strings.Contains(r.Header.Get("Accept"), "application/json")

	if isAJAX {
		helpers.SendJson(w, &helpers.Response{
			Status:  http.StatusUnauthorized,
			Message: "Invalid or expired session",
		})
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
