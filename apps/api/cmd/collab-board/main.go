package main

import (
	"log"

	db "github.com/mirkosisko-dev/api/db"
	"github.com/mirkosisko-dev/api/internal/api"
)

func main() {
	db, err := db.NewPostgreSQLStorage()
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(":8080", db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
