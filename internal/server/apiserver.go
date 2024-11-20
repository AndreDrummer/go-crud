package server

import (
	"fmt"
	"go-crud/internal/server/app"
)

func RunServer() error {
	httpServer := app.New()

	if err := httpServer.ListenAndServe(); err != nil {
		return fmt.Errorf("an error ocurred when starting the server: %v", err)
	}

	return nil
}
