package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	db "github.com/mirkosisko-dev/api/db"
	"github.com/mirkosisko-dev/api/internal/api/service/organization"
	"github.com/mirkosisko-dev/api/internal/api/service/user"
)

type APIServer struct {
	addr string
	db   *db.Database
}

func NewAPIServer(addr string, db *db.Database) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userHandler := user.NewHandler()
	userHandler.RegisterRoutes(subrouter)

	organizationHandler := organization.NewHandler()
	organizationHandler.RegisterRoutes(subrouter)

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, subrouter)
}
