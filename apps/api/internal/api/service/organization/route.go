package organization

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mirkosisko-dev/api/db"
	"github.com/mirkosisko-dev/api/db/sqlc"
	"github.com/mirkosisko-dev/api/internal/api/service/auth"
	"github.com/mirkosisko-dev/api/internal/types"
	"github.com/mirkosisko-dev/api/utils"
)

type Handler struct {
	storage *pool.Database
}

func NewHandler(storage *pool.Database) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/organization", h.handleCreateOrganization).Methods(http.MethodPost)
	router.HandleFunc("/organization/{orgId}/invites", h.handleCreateInvite).Methods(http.MethodPost)
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

func (h *Handler) handleCreateInvite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orgID := vars["orgId"]

	if orgID == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("organization id is required"))
		return
	}

	orgIdInt, err := strconv.ParseInt(orgID, 10, 32)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("invalid organization id"))
		return
	}

	var payload types.CreateInvitePayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	invite, err := h.storage.Query.CreateOrganizationInvite(r.Context(), sqlc.CreateOrganizationInviteParams{
		OrganizationID:  pgtype.Int4{Int32: int32(orgIdInt), Valid: true},
		InvitedUserID:   pgtype.Int4{Int32: payload.InvitedUserID, Valid: true},
		InvitedByUserID: pgtype.Int4{Int32: int32(userID), Valid: true},
		ExpiresAt:       pgtype.Timestamp{Time: payload.ExpiresAt, Valid: true},
		Role:            sqlc.OrganizationRoleOwner,
		Status:          sqlc.OrganizationInviteStatusPending,
	})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, invite)
}
