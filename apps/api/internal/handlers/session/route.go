package session

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgtype"
	pool "github.com/mirkosisko-dev/api/db"
	"github.com/mirkosisko-dev/api/db/sqlc"
	"github.com/mirkosisko-dev/api/internal/handlers/auth"
	"github.com/mirkosisko-dev/api/utils"
)

type Handler struct {
	storage *pool.Database
}

func NewHandler(storage *pool.Database) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/session", h.handleCreateSession).Methods("POST")
	router.HandleFunc("/session", h.handleGetSession).Methods("GET")
	router.HandleFunc("/session/revoke", h.handleRevokeSession).Methods("POST")
	router.HandleFunc("/session/{sessionID}", h.handleDeleteSession).Methods("DELETE")
}

func (h *Handler) handleCreateSession(w http.ResponseWriter, r *http.Request) {
	var payload sqlc.CreateSesionParams
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	_, err := h.storage.Query.CreateSesion(r.Context(), sqlc.CreateSesionParams{
		RefreshToken: payload.RefreshToken,
		IsRevoked:    payload.IsRevoked,
		ExpiresAt:    payload.ExpiresAt,
	})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleGetSession(w http.ResponseWriter, r *http.Request) {
	type getSessionPayload struct {
		ID uuid.UUID `json:"uuid"`
	}

	var payload getSessionPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	session, err := h.storage.Query.GetSession(r.Context(), pgtype.UUID{Bytes: payload.ID, Valid: true})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, session)
}

func (h *Handler) handleRevokeSession(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(auth.AuthKey{}).(*auth.TokenClaims)

	sessionUUID, err := uuid.Parse(claims.RegisteredClaims.ID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	err = h.storage.Query.RevokeSession(r.Context(), pgtype.UUID{Bytes: sessionUUID, Valid: true})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleDeleteSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionIDStr := vars["sessionID"]
	if sessionIDStr == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("session id is required"))
		return
	}

	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("invalid session id"))
		return
	}

	err = h.storage.Query.DeleteSession(r.Context(), pgtype.UUID{Bytes: sessionID, Valid: true})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}
