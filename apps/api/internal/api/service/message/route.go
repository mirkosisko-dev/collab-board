package message

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mirkosisko-dev/api/db"
	"github.com/mirkosisko-dev/api/db/sqlc"
	"github.com/mirkosisko-dev/api/utils"
)

type Handler struct {
	storage *pool.Database
}

func NewHandler(storage *pool.Database) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/message", h.handleCreateMessage).Methods("POST")
}

func (h *Handler) handleCreateMessage(w http.ResponseWriter, r *http.Request) {
	var payload sqlc.CreateMessageParams

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	_, err := h.storage.Query.CreateMessage(r.Context(), sqlc.CreateMessageParams{
		BoardID:        payload.BoardID,
		OrganizationID: payload.OrganizationID,
		UserID:         payload.UserID,
		DocumentID:     payload.DocumentID,
		Content:        payload.Content,
	})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}
