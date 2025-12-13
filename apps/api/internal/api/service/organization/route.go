package organization

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mirkosisko-dev/api/db"
	"github.com/mirkosisko-dev/api/utils"
)

type Handler struct {
	storage *pool.Database
}

func NewHandler(storage *pool.Database) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/organization", h.handleCreateOrganization).Methods("POST")
}

func (h *Handler) handleCreateOrganization(w http.ResponseWriter, r *http.Request) {
	type createOrgPayload struct {
		Name string `json:"name"`
	}

	var payload createOrgPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	_, err := h.storage.Query.CreateOrganization(r.Context(), payload.Name)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleInviteToOrganization(w http.ResponseWriter, r *http.Request) {
}
