package app

import (
	"go-crud/internal/server/app/repository"
	userrepository "go-crud/internal/server/app/repository/User"
	"go-crud/internal/server/app/router"
	"net/http"
	"time"
)

func New() *http.Server {

	serverHandler := router.Handler(
		map[string]repository.Repository{
			repository.User: userrepository.NewUserRepository(),
		},
	)

	return &http.Server{
		Addr:         ":8080",
		Handler:      serverHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Minute,
	}
}
