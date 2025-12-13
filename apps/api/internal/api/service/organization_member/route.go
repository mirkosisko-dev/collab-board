package organizationmember

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
	router.HandleFunc("/organization-member", h.handleCreateOrganizationMember).Methods("POST")
}

func (h *Handler) handleCreateOrganizationMember(w http.ResponseWriter, r *http.Request) {
	var payload sqlc.CreateOrganizationMemberParams
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	_, err := h.storage.Query.CreateOrganizationMember(r.Context(), sqlc.CreateOrganizationMemberParams{
		OrganizationID: payload.OrganizationID,
		UserID:         payload.UserID,
		Role:           payload.Role,
	})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}
