package router

import (
	"encoding/json"
	"errors"
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
	r.Get("/users", FindAll())
	r.Get("/users/{id}", FindByID())
	r.Post("/users", Insert())
	r.Put("/users/{id}", Update())
	r.Delete("/users/{id}", Delete())
}

func handleBodyRequest(w http.ResponseWriter, r *http.Request, user *model.User) error {
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		msgError := fmt.Sprintf("Invalid JSON: %v", err)

		slog.Error(msgError)
		sendReponse(w, Response{Error: "Invalid request: body malformed"}, http.StatusBadRequest)
		return errors.New(msgError)
	}

	if !user.IsValid() {
		msgError := "user Invalid: Missing required fields"

		slog.Error(msgError)
		sendReponse(w, Response{Error: "Please provide first name, last name and biography for the user"}, http.StatusBadRequest)
		return errors.New(msgError)
	}

	return nil
}

func FindAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		db := dbhandler.OpenDB()
		data, err := db.FindAll()

		if err != nil {
			slog.Error("ERROR", "", err)
			sendReponse(w, Response{Error: "The users information could not be retrieved"}, http.StatusInternalServerError)
			return
		}

		users := make([]model.User, len(data))
		for i, v := range data {
			if err := json.Unmarshal([]byte(v), &users[i]); err != nil {
				slog.Error("error converting the data from DB to JSON", "error", err)
				sendReponse(w, Response{Error: "The users information could not be retrieved"}, http.StatusInternalServerError)
				return
			}
		}

		slog.Info("SUCESS", "Users", users)
		sendReponse(w, Response{Data: users}, http.StatusOK)
	}
}

func FindByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		userID := chi.URLParam(r, "id")

		db := dbhandler.OpenDB()
		userString, err := db.FindByID(userID)

		if err != nil {
			if errors.Is(err, &dbhandler.DBNotFoundError{}) {
				slog.Error("User was not found")
				sendReponse(w, Response{Error: "The user with the specified ID does not exist"}, http.StatusNotFound)
			} else {
				slog.Error(fmt.Sprintf("an error occurred tryna get the user %v", userID), "error", err)
				sendReponse(w, Response{Error: "The user information could not be retrieved"}, http.StatusInternalServerError)
			}
			return
		}

		var user model.User
		if err := json.Unmarshal([]byte(userString), &user); err != nil {
			slog.Error("error converting the response to JSON", "error", err)
			sendReponse(w, Response{Error: "The user information could not be retrieved"}, http.StatusInternalServerError)
			return
		}

		slog.Info("SUCCES", "User found", user)
		sendReponse(w, Response{Data: user}, http.StatusOK)

	}
}

func Insert() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var user model.User

		if err := handleBodyRequest(w, r, &user); err == nil {
			intUUID := uuid.New()
			userID := intUUID.String()
			user.ID = userID

			userJson, err := json.Marshal(user)

			if err != nil {
				sendReponse(w, Response{Error: fmt.Sprintf("database error: %v", err)}, http.StatusInternalServerError)
				slog.Error("error parsing json user", "error", err)
				return
			}

			db := dbhandler.OpenDB()

			if err := db.Insert(string(userJson)); err != nil {
				slog.Error("error inserting user in DB", "error", err)
				sendReponse(w, Response{Error: "There was an error while saving the user to the database"}, http.StatusInternalServerError)
				return
			}

			slog.Info("SUCCES", "User created succesfully", user)
			sendReponse(w, Response{Data: user}, http.StatusCreated)
		}
	}
}

func Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var user model.User

		if err := handleBodyRequest(w, r, &user); err == nil {
			userID := chi.URLParam(r, "id")
			user.ID = userID

			userJson, err := json.Marshal(user)
			if err != nil {
				sendReponse(w, Response{Error: fmt.Sprintf("database error: %v", err)}, http.StatusInternalServerError)
				slog.Error("error parsing json user", "error", err)
				return
			}

			db := dbhandler.OpenDB()

			if err := db.Update(userID, string(userJson)); err != nil {
				if errors.Is(err, &dbhandler.DBNotFoundError{}) {
					slog.Error("User was not found")
					sendReponse(w, Response{Error: "The user with the specified ID does not exist"}, http.StatusNotFound)
				} else {
					slog.Error("error inserting user in DB", "error", err)
					sendReponse(w, Response{Error: "The user information could not be modified"}, http.StatusInternalServerError)
				}
				return
			} else {
				slog.Info("SUCCES", "User updated", user)
				sendReponse(w, Response{Data: user}, http.StatusOK)
			}
		}
	}
}

func Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		userID := chi.URLParam(r, "id")

		db := dbhandler.OpenDB()
		err := db.Delete(userID)

		if err != nil {
			if errors.Is(err, &dbhandler.DBNotFoundError{}) {
				slog.Error("User was not found")
				sendReponse(w, Response{Error: "The user with the specified ID does not exist"}, http.StatusNotFound)
			} else {
				slog.Error("Operation error", "error", err)
				sendReponse(w, Response{Error: "The user could not be removed"}, http.StatusInternalServerError)
			}
			return
		}

		slog.Info("SUCCES", "User removed", nil)
	}
}
