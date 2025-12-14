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
	"github.com/mirkosisko-dev/api/internal/middleware"
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
	subrouter := router.PathPrefix("/v1").Subrouter()

	public := subrouter.NewRoute().Subrouter()
	protected := subrouter.NewRoute().Subrouter()

	protected.Use(middleware.AuthenticationMiddleware)

	boardHandler := board.NewHandler(s.db)
	boardHandler.RegisterRoutes(protected)

	boardColumnHandler := boardcolumn.NewHandler(s.db)
	boardColumnHandler.RegisterRoutes(protected)

	documentHandler := document.NewHandler(s.db)
	documentHandler.RegisterRoutes(protected)

	documentContentHandler := documentcontent.NewHandler(s.db)
	documentContentHandler.RegisterRoutes(protected)

	messageHandler := message.NewHandler(s.db)
	messageHandler.RegisterRoutes(protected)

	organizationHandler := organization.NewHandler(s.db)
	organizationHandler.RegisterRoutes(protected)

	organizationMemberHandler := organizationmember.NewHandler(s.db)
	organizationMemberHandler.RegisterRoutes(protected)

	taskHandler := task.NewHandler(s.db)
	taskHandler.RegisterRoutes(protected)

	userHandler := user.NewHandler(s.db)
	userHandler.RegisterPublicRoutes(public)
	userHandler.RegisterProtectedRoutes(protected)

	log.Println("Listening on", s.addr)

	handler := middleware.CORS(router)

	return http.ListenAndServe(s.addr, handler)
}
