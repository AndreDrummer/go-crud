package main

import (
	"fmt"
	"go-crud/go-crud/server"
	"net/http"
	"time"
)

func main() {
	if err := run(); err != nil {
		return
	}
}

func run() error {
	serverHandler := server.Handler()

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
