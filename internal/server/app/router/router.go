package router

import (
	"encoding/json"
	"errors"
	"fmt"
	customerrors "go-crud/internal/server/app/errors"
	"go-crud/internal/server/app/model"
	"go-crud/internal/server/app/repository"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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

func userOperations(r chi.Router, repository repository.Repository) {
	r.Get("/users", FindAll(repository))
	r.Get("/users/{id}", FindByID(repository))
	r.Post("/users", Insert(repository))
	r.Put("/users/{id}", Update(repository))
	r.Delete("/users/{id}", Delete(repository))
}

func Handler(repositories map[string]repository.Repository) http.Handler {
	handler := chi.NewMux()

	handler.Use(middleware.Recoverer)
	handler.Use(middleware.RequestID)
	handler.Use(middleware.Logger)

	handler.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			userOperations(r, repositories[repository.User])
		})

	})

	return handler
}

func FindAll(repository repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		userList, err := repository.GetAll()

		if err != nil {
			slog.Error(err.Error())

			sendReponse(
				w,
				Response{Error: "The users information could not be retrieved"},
				http.StatusInternalServerError,
			)
			return
		}

		slog.Info(fmt.Sprintf("Users %v", userList))
		sendReponse(w, Response{Data: userList}, http.StatusOK)
	}
}

func FindByID(repository repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		userID := chi.URLParam(r, "id")

		user, err := repository.GetOne(userID)

		if err != nil {
			NotFoundError := &customerrors.NotFoundError{}

			if errors.As(err, &NotFoundError) {
				slog.Error("User was not found")

				sendReponse(
					w,
					Response{Error: "The user with the specified ID does not exist"},
					http.StatusNotFound,
				)
			} else {
				slog.Error(err.Error())

				sendReponse(
					w,
					Response{Error: "The user information could not be retrieved"},
					http.StatusInternalServerError,
				)
			}
			return
		}

		slog.Info(fmt.Sprintf("User %v", user))
		sendReponse(w, Response{Data: user}, http.StatusOK)

	}
}

func Insert(repository repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var user model.User

		if err := handleBodyRequest(w, r, &user); err == nil {
			err := repository.Insert(&user)

			if err != nil {
				slog.Error(err.Error())
				sendReponse(
					w,
					Response{Error: "There was an error while saving the user to the database"},
					http.StatusInternalServerError,
				)
				return
			}

			slog.Info(fmt.Sprintf("User created succesfully %v", user))
			sendReponse(w, Response{Data: user}, http.StatusCreated)
		}
	}
}

func Update(repository repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var user model.User

		if err := handleBodyRequest(w, r, &user); err == nil {
			userID := chi.URLParam(r, "id")
			user.ID = userID

			err := repository.Update(&user)

			if err != nil {
				NotFoundError := &customerrors.NotFoundError{}
				if errors.As(err, &NotFoundError) {
					slog.Error("User was not found")
					sendReponse(w, Response{Error: "The user with the specified ID does not exist"}, http.StatusNotFound)
				} else {
					slog.Error(err.Error())
					sendReponse(w, Response{Error: "The user information could not be modified"}, http.StatusInternalServerError)
				}
				return
			} else {
				slog.Info(fmt.Sprintf("User updated %v", user))
				sendReponse(w, Response{Data: user}, http.StatusOK)
			}
		}
	}
}

func Delete(repository repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		userID := chi.URLParam(r, "id")

		err := repository.Delete(userID)

		if err != nil {
			NotFoundError := &customerrors.NotFoundError{}
			if errors.As(err, &NotFoundError) {
				slog.Error("User was not found")
				sendReponse(w, Response{Error: "The user with the specified ID does not exist"}, http.StatusNotFound)
			} else {
				slog.Error(err.Error())
				sendReponse(w, Response{Error: "The user could not be removed"}, http.StatusInternalServerError)
			}
			return
		}

		slog.Info("User removed successfully")
	}
}
