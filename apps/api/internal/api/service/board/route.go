package board

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
	router.HandleFunc("/board", h.handleCreateBoard).Methods("POST")
}

func (h *Handler) handleCreateBoard(w http.ResponseWriter, r *http.Request) {

	var payload sqlc.CreateBoardParams
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	_, err := h.storage.Query.CreateBoard(r.Context(), sqlc.CreateBoardParams{
		OrganizationID: payload.OrganizationID,
		Name:           payload.Name,
		CreatedBy:      payload.CreatedBy,
	})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}
