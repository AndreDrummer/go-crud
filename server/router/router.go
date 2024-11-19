package router

import (
	"encoding/json"
	"fmt"
	dbhandler "go-crud/server/db/handler"
	"go-crud/server/model"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
)

type Response struct {
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
}

func sendReponse(w http.ResponseWriter, resp Response, status int) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(resp)

	if err != nil {
		slog.Error("error parsing response", "error", err)
		sendReponse(
			w,
			Response{Error: "something went wrong!"},
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		slog.Error("error writing response", "error", err)
		return
	}
}

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
	r.Patch("/users/{id}", updateUser())
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
			slog.Error(fmt.Sprintf("Invalid JSON: %v", err))
			sendReponse(w, Response{Error: "Invalid request: body malformed"}, http.StatusBadRequest)
			return
		}

		if !user.IsValid() {
			slog.Error("User Invalid: Missing required fields.")
			sendReponse(w, Response{Error: "Please provide first name, last name and biography"}, http.StatusBadRequest)
			return
		} else {
			intUUID := uuid.New()
			userID := intUUID.String()
			user.ID = userID

			userStringLine := dbhandler.AnyToString(user)
			db := dbhandler.OpenDB()

			if err := db.Insert(userStringLine); err != nil {
				sendReponse(w, Response{Error: fmt.Sprintf("database error: %v", err)}, http.StatusBadRequest)
				return
			}

			sendReponse(w, Response{Data: "User created successfully"}, http.StatusCreated)
			slog.Info("User created successfully")
		}

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
