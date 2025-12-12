package organization

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mirkosisko-dev/api/db"
	"github.com/mirkosisko-dev/api/utils"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/organization", h.handleCreateOrganization).Methods("POST")
}

func (h *Handler) handleCreateOrganization(w http.ResponseWriter, r *http.Request) {
	storage, err := pool.NewPostgreSQLStorage()
	if err != nil {
		log.Fatal(err)
	}

	type createOrgPayload struct {
		Name string `json:"name"`
	}

	var payload createOrgPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	_, err = storage.Query.CreateOrganization(context.Background(), payload.Name)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleInviteToOrganization(w http.ResponseWriter, r *http.Request) {
}
