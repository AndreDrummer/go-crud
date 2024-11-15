package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

func Handler() http.Handler {
	handler := chi.NewMux()

	handler.Use(middleware.Recoverer)
	handler.Use(middleware.RequestID)
	handler.Use(middleware.Logger)

	handler.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Get("/users", readUsers())
			r.Get("/user/{id[0-9]+}", getUser())
			r.Post("/users", createUser())
			r.Patch("/users", updateUser())
			r.Delete("/users", deleteUser())
		})

	})

	return handler

}

func getUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func createUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func readUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Works!")
	}
}

func updateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func deleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
