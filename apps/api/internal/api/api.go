package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mirkosisko-dev/api/db"
	"github.com/mirkosisko-dev/api/internal/api/service/board"
	boardcolumn "github.com/mirkosisko-dev/api/internal/api/service/board_column"
	"github.com/mirkosisko-dev/api/internal/api/service/document"
	documentcontent "github.com/mirkosisko-dev/api/internal/api/service/document_content"
	"github.com/mirkosisko-dev/api/internal/api/service/message"
	"github.com/mirkosisko-dev/api/internal/api/service/organization"
	organizationmember "github.com/mirkosisko-dev/api/internal/api/service/organization_member"
	"github.com/mirkosisko-dev/api/internal/api/service/task"
	"github.com/mirkosisko-dev/api/internal/api/service/user"
)

type APIServer struct {
	addr string
	db   *pool.Database
}

func NewAPIServer(addr string, db *pool.Database) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	boardHandler := board.NewHandler()
	boardHandler.RegisterRoutes(subrouter)

	boardColumnHandler := boardcolumn.NewHandler()
	boardColumnHandler.RegisterRoutes(subrouter)

	documentHandler := document.NewHandler()
	documentHandler.RegisterRoutes(subrouter)

	documentContentHandler := documentcontent.NewHandler()
	documentContentHandler.RegisterRoutes(subrouter)

	messageHandler := message.NewHandler()
	messageHandler.RegisterRoutes(subrouter)

	organizationHandler := organization.NewHandler()
	organizationHandler.RegisterRoutes(subrouter)

	organizationMemberHandler := organizationmember.NewHandler()
	organizationMemberHandler.RegisterRoutes(subrouter)

	taskHandler := task.NewHandler()
	taskHandler.RegisterRoutes(subrouter)

	userHandler := user.NewHandler()
	userHandler.RegisterRoutes(subrouter)

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, subrouter)
}
