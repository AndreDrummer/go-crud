package server

import (
	"fmt"
	"go-crud/server/router"
	"net/http"
	"time"
)

func RunServer() error {
	serverHandler := router.Handler()

	server := &http.Server{
		Addr:         ":8080",
		Handler:      serverHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	if err := server.ListenAndServe(); err != nil {
		return fmt.Errorf("an error ocurred when starting the server: %v", err)
	}

	return nil
}
