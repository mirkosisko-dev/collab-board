package api

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/mirkosisko-dev/api/config"
	pool "github.com/mirkosisko-dev/api/db"
	"github.com/mirkosisko-dev/api/internal/handlers/board"
	boardcolumn "github.com/mirkosisko-dev/api/internal/handlers/board_column"
	"github.com/mirkosisko-dev/api/internal/handlers/document"
	documentcontent "github.com/mirkosisko-dev/api/internal/handlers/document_content"
	"github.com/mirkosisko-dev/api/internal/handlers/message"
	"github.com/mirkosisko-dev/api/internal/handlers/organization"
	organizationmember "github.com/mirkosisko-dev/api/internal/handlers/organization_member"
	"github.com/mirkosisko-dev/api/internal/handlers/session"
	"github.com/mirkosisko-dev/api/internal/handlers/task"
	"github.com/mirkosisko-dev/api/internal/handlers/user"
	"github.com/mirkosisko-dev/api/internal/middleware"
	sessionservice "github.com/mirkosisko-dev/api/internal/services/session"
)

type APIServer struct {
	addr   string
	db     *pool.Database
	config *config.Config
}

func NewAPIServer(addr string, db *pool.Database, cfg *config.Config) *APIServer {
	return &APIServer{
		addr:   addr,
		db:     db,
		config: cfg,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/v1").Subrouter()

	subrouter.Use(middleware.Logger)

	public := subrouter.NewRoute().Subrouter()
	protected := subrouter.NewRoute().Subrouter()

	protected.Use(middleware.AuthenticationMiddleware(s.config.AccessTokenSecret))

	boardhandlers := board.NewHandler(s.db)
	boardhandlers.RegisterRoutes(protected)

	boardColumnhandlers := boardcolumn.NewHandler(s.db)
	boardColumnhandlers.RegisterRoutes(protected)

	documenthandlers := document.NewHandler(s.db)
	documenthandlers.RegisterRoutes(protected)

	documentContenthandlers := documentcontent.NewHandler(s.db)
	documentContenthandlers.RegisterRoutes(protected)

	messagehandlers := message.NewHandler(s.db)
	messagehandlers.RegisterRoutes(protected)

	organizationhandlers := organization.NewHandler(s.db)
	organizationhandlers.RegisterRoutes(protected)

	organizationMemberhandlers := organizationmember.NewHandler(s.db)
	organizationMemberhandlers.RegisterRoutes(protected)

	taskHandler := task.NewHandler(s.db)
	taskHandler.RegisterRoutes(protected)

	sessionService := sessionservice.NewService(s.db, s.config)

	userHandler := user.NewHandler(s.db, s.config, sessionService)
	userHandler.RegisterPublicRoutes(public)
	userHandler.RegisterProtectedRoutes(protected)

	sessionshandlers := session.NewHandler(s.db)
	sessionshandlers.RegisterRoutes(protected)

	handlers := middleware.CORS(router)

	srv := &http.Server{
		Addr:    s.addr,
		Handler: handlers,
	}

	go func() {
		slog.Info("server started", "addr", s.addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("listen", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		return err
	}

	slog.Info("server exiting")
	return nil
}
