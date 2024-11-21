package router

import (
	"go-crud/internal/server/app/repository"
	routes "go-crud/internal/server/app/router/Routes"

	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func Handler(repositories map[string]repository.Repository) http.Handler {
	handler := chi.NewMux()

	handler.Use(middleware.Recoverer)
	handler.Use(middleware.RequestID)
	handler.Use(middleware.Logger)

	handler.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			routes.UserRoutes(r, repositories[repository.User])
		})

	})

	return handler
}
