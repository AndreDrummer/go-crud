package router_http

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-crud/internal/server/app/model"
	"log/slog"
	"net/http"
)

type Response struct {
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
}

func Send(w http.ResponseWriter, resp Response, status int) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(resp)

	if err != nil {
		slog.Error("error parsing response", "error", err)
		Send(
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
		Send(w, Response{Error: "Invalid request: body malformed"}, http.StatusBadRequest)
		return errors.New(msgError)
	}

	if !user.IsValid() {
		msgError := "user Invalid: Missing required fields"

		slog.Error(msgError)
		Send(w, Response{Error: "Please provide first name, last name and biography for the user"}, http.StatusBadRequest)
		return errors.New(msgError)
	}

	return nil
}
