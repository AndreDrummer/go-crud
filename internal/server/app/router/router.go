package router

import (
	"encoding/json"
	"errors"
	"fmt"
	dberrors "go-crud/internal/server/app/db/errors"
	dbhandler "go-crud/internal/server/app/db/handler"
	"go-crud/internal/server/app/model"
	"go-crud/internal/server/app/repository"
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

func SendReponse(w http.ResponseWriter, resp Response, status int) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(resp)

	if err != nil {
		slog.Error("error parsing response", "error", err)
		SendReponse(
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

func HandleBodyRequest(w http.ResponseWriter, r *http.Request, user *model.User) error {
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		msgError := fmt.Sprintf("Invalid JSON: %v", err)

		slog.Error(msgError)
		SendReponse(w, Response{Error: "Invalid request: body malformed"}, http.StatusBadRequest)
		return errors.New(msgError)
	}

	if !user.IsValid() {
		msgError := "user Invalid: Missing required fields"

		slog.Error(msgError)
		SendReponse(w, Response{Error: "Please provide first name, last name and biography for the user"}, http.StatusBadRequest)
		return errors.New(msgError)
	}

	return nil
}

func userOperations(r chi.Router) {
	r.Get("/users", FindAll())
	r.Get("/users/{id}", FindByID())
	r.Post("/users", Insert())
	r.Put("/users/{id}", Update())
	r.Delete("/users/{id}", Delete())
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

func FindAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		userList, err := repository.GetUserList()

		if err != nil {
			slog.Error(err.Error())

			SendReponse(
				w,
				Response{Error: "The users information could not be retrieved"},
				http.StatusInternalServerError,
			)
			return
		}

		slog.Info(fmt.Sprintf("Users %v", userList))
		SendReponse(w, Response{Data: userList}, http.StatusOK)
	}
}

func FindByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		userID := chi.URLParam(r, "id")

		user, err := repository.GetUser(userID)

		if err != nil {
			DBNotFoundError := &dberrors.DBNotFoundError{}

			if errors.As(err, &DBNotFoundError) {
				slog.Error("User was not found")

				SendReponse(
					w,
					Response{Error: "The user with the specified ID does not exist"},
					http.StatusNotFound,
				)
			} else {
				slog.Error(err.Error())

				SendReponse(
					w,
					Response{Error: "The users information could not be retrieved"},
					http.StatusInternalServerError,
				)
			}
			return
		}

		slog.Info(fmt.Sprintf("User %v", user))
		SendReponse(w, Response{Data: user}, http.StatusOK)

	}
}

func Insert() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var user model.User

		if err := HandleBodyRequest(w, r, &user); err == nil {
			intUUID := uuid.New()
			userID := intUUID.String()
			user.ID = userID

			userJson, err := json.Marshal(user)

			if err != nil {
				SendReponse(w, Response{Error: fmt.Sprintf("database error: %v", err)}, http.StatusInternalServerError)
				slog.Error("error parsing json user", "error", err)
				return
			}

			db := dbhandler.OpenDB()

			if err := db.Insert(string(userJson)); err != nil {
				slog.Error("error inserting user in DB", "error", err)
				SendReponse(w, Response{Error: "There was an error while saving the user to the database"}, http.StatusInternalServerError)
				return
			}

			slog.Info("SUCCES", "User created succesfully", user)
			SendReponse(w, Response{Data: user}, http.StatusCreated)
		}
	}
}

func Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var user model.User

		if err := HandleBodyRequest(w, r, &user); err == nil {
			userID := chi.URLParam(r, "id")
			user.ID = userID

			userJson, err := json.Marshal(user)
			if err != nil {
				SendReponse(w, Response{Error: fmt.Sprintf("database error: %v", err)}, http.StatusInternalServerError)
				slog.Error("error parsing json user", "error", err)
				return
			}

			db := dbhandler.OpenDB()

			if err := db.Update(userID, string(userJson)); err != nil {
				if errors.Is(err, &dberrors.DBNotFoundError{}) {
					slog.Error("User was not found")
					SendReponse(w, Response{Error: "The user with the specified ID does not exist"}, http.StatusNotFound)
				} else {
					slog.Error("error inserting user in DB", "error", err)
					SendReponse(w, Response{Error: "The user information could not be modified"}, http.StatusInternalServerError)
				}
				return
			} else {
				slog.Info("SUCCES", "User updated", user)
				SendReponse(w, Response{Data: user}, http.StatusOK)
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
			if errors.Is(err, &dberrors.DBNotFoundError{}) {
				slog.Error("User was not found")
				SendReponse(w, Response{Error: "The user with the specified ID does not exist"}, http.StatusNotFound)
			} else {
				slog.Error("Operation error", "error", err)
				SendReponse(w, Response{Error: "The user could not be removed"}, http.StatusInternalServerError)
			}
			return
		}

		slog.Info("SUCCES", "User removed", nil)
	}
}
