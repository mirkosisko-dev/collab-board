package main

import (
	"log"

	db "github.com/mirkosisko-dev/api/db"
)

func main() {
	_, err := db.NewPostgreSQLStorage()
	if err != nil {
		log.Fatal(err)
	}
}
