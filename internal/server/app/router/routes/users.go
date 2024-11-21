package routes

import (
	"errors"
	"fmt"
	customerrors "go-crud/internal/server/app/errors"
	"go-crud/internal/server/app/model"
	"go-crud/internal/server/app/repository"
	router_http "go-crud/internal/server/app/router/http"

	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
)

func UserRoutes(r chi.Router, repository repository.Repository) {
	r.Get("/users", FindAll(repository))
	r.Get("/users/{id}", FindByID(repository))
	r.Post("/users", Insert(repository))
	r.Put("/users/{id}", Update(repository))
	r.Delete("/users/{id}", Delete(repository))
}

func FindAll(repository repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		userList, err := repository.GetAll()

		if err != nil {
			slog.Error(err.Error())

			router_http.Send(
				w,
				router_http.Response{Error: "The users information could not be retrieved"},
				http.StatusInternalServerError,
			)
			return
		}

		slog.Info(fmt.Sprintf("Users %v", userList))
		router_http.Send(w, router_http.Response{Data: userList}, http.StatusOK)
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

				router_http.Send(
					w,
					router_http.Response{Error: "The user with the specified ID does not exist"},
					http.StatusNotFound,
				)
			} else {
				slog.Error(err.Error())

				router_http.Send(
					w,
					router_http.Response{Error: "The user information could not be retrieved"},
					http.StatusInternalServerError,
				)
			}
			return
		}

		slog.Info(fmt.Sprintf("User %v", user))
		router_http.Send(w, router_http.Response{Data: user}, http.StatusOK)

	}
}

func Insert(repository repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var user model.User

		if err := router_http.HandleBodyRequest(w, r, &user); err == nil {
			err := repository.Insert(&user)

			if err != nil {
				slog.Error(err.Error())
				router_http.Send(
					w,
					router_http.Response{Error: "There was an error while saving the user to the database"},
					http.StatusInternalServerError,
				)
				return
			}

			slog.Info(fmt.Sprintf("User created succesfully %v", user))
			router_http.Send(w, router_http.Response{Data: user}, http.StatusCreated)
		}
	}
}

func Update(repository repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var user model.User

		if err := router_http.HandleBodyRequest(w, r, &user); err == nil {
			userID := chi.URLParam(r, "id")
			user.ID = userID

			err := repository.Update(&user)

			if err != nil {
				NotFoundError := &customerrors.NotFoundError{}
				if errors.As(err, &NotFoundError) {
					slog.Error("User was not found")
					router_http.Send(w, router_http.Response{Error: "The user with the specified ID does not exist"}, http.StatusNotFound)
				} else {
					slog.Error(err.Error())
					router_http.Send(w, router_http.Response{Error: "The user information could not be modified"}, http.StatusInternalServerError)
				}
				return
			} else {
				slog.Info(fmt.Sprintf("User updated %v", user))
				router_http.Send(w, router_http.Response{Data: user}, http.StatusOK)
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
				router_http.Send(w, router_http.Response{Error: "The user with the specified ID does not exist"}, http.StatusNotFound)
			} else {
				slog.Error(err.Error())
				router_http.Send(w, router_http.Response{Error: "The user could not be removed"}, http.StatusInternalServerError)
			}
			return
		}

		slog.Info("User removed successfully")
	}
}
