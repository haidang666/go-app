package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/haidang666/go-app/internal/infrastructure/http/module_route/auth"
)

type NewRouterArgs struct {
	AuthHandler        *auth.AuthHandler
}

func NewRouter(args NewRouterArgs) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	})

	r.Route("/api/v1", func(ur chi.Router) {
		auth.RegisterRoutes(ur, args.AuthHandler)
	})

	return r
}
