package main

import (
	"go-crud/internal/server"
	"log/slog"
)

func main() {
	if err := server.RunServer(); err != nil {
		slog.Error("an error running server", "error", err)
		return
	}
}
