package task

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
	router.HandleFunc("/task", h.handleCreateTask).Methods("POST")
}

func (h *Handler) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var payload sqlc.CreateTaskParams
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	_, err := h.storage.Query.CreateTask(r.Context(), sqlc.CreateTaskParams{
		BoardID:     payload.BoardID,
		ColumnID:    payload.ColumnID,
		AssigneeID:  payload.AssigneeID,
		CreatedBy:   payload.CreatedBy,
		Title:       payload.Title,
		Description: payload.Description,
		Position:    payload.Position,
	})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}
