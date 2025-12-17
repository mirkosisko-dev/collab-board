package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/mirkosisko-dev/api/config"
	db "github.com/mirkosisko-dev/api/db"
	"github.com/mirkosisko-dev/api/internal/api"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg := config.Load()

	db, err := db.NewPostgreSQLStorage()
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(":8080", db, &cfg)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
