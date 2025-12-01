package auth

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, h *AuthHandler) {
	r.Route("/auth", func(ur chi.Router) {
		ur.Post("/sign-up", h.SignUp)
	})
}
