package router

import (
	"encoding/json"
	"fmt"
	dbhandler "go-crud/server/db/handler"
	"go-crud/server/model"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
)

func Handler() http.Handler {
	handler := chi.NewMux()

	handler.Use(middleware.Recoverer)
	handler.Use(middleware.RequestID)
	handler.Use(middleware.Logger)

	handler.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			userOperations(r)
		})

	})

	return handler

}
func userOperations(r chi.Router) {
	r.Get("/users", readUsers())
	r.Get("/users/{id}", getUser())
	r.Post("/users", createUser())
	r.Patch("/users", updateUser())
	r.Delete("/users/{id}", deleteUser())
}

func getUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
	}
}

func createUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var user model.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
			return
		}

		intUUID := uuid.New()
		userID := intUUID.String()
		user.ID = userID

		userStringLine := dbhandler.AnyToString(user)
		db := dbhandler.OpenDB()
		db.Insert(userStringLine)
	}
}

func readUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		fmt.Fprint(w, "Works!")
	}
}

func updateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
	}
}

func deleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
	}
}
