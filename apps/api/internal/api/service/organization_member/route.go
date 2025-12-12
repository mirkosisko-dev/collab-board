package organizationmember

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mirkosisko-dev/api/db"
	"github.com/mirkosisko-dev/api/db/sqlc"
	"github.com/mirkosisko-dev/api/utils"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/organization-member", h.handleCreateOrganizationMember).Methods("POST")
}

func (h *Handler) handleCreateOrganizationMember(w http.ResponseWriter, r *http.Request) {
	storage, err := pool.NewPostgreSQLStorage()
	if err != nil {
		log.Fatal(err)
	}

	var payload sqlc.CreateOrganizationMemberParams
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	_, err = storage.Query.CreateOrganizationMember(context.Background(), sqlc.CreateOrganizationMemberParams{
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
