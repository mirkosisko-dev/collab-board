package user

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
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

func (h *Handler) RegisterPublicRoutes(r *mux.Router) {
	r.HandleFunc("/auth/register", h.handleRegister).Methods(http.MethodPost)
	r.HandleFunc("/auth/login", h.handleLogin).Methods(http.MethodPost)
}

func (h *Handler) RegisterProtectedRoutes(r *mux.Router) {
	// r.HandleFunc("/users/me", h.handleMe).Methods(http.MethodGet)
	// r.HandleFunc("/users/me/password", h.handleUpdatePassword).Methods(http.MethodPatch)
	r.HandleFunc("/invites", h.handleGetInvites).Methods(http.MethodGet)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", payload))
		return
	}

	user, err := h.storage.Query.GetUserByEmail(r.Context(), payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("not found, invalid email or password"))
		return
	}

	if !auth.ComparePasswords(user.PasswordHash, []byte(payload.Password)) {
		utils.WriteError(w, http.StatusBadRequest, errors.New("not found, invalid email or password"))
		return
	}

	userUUID, err := uuid.Parse(user.ID.String())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}

	token, err := auth.CreateJWT(userUUID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterUserPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	_, err := h.storage.Query.GetUserByEmail(r.Context(), payload.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	_, err = h.storage.Query.CreateUser(r.Context(), sqlc.CreateUserParams{
		Name:         payload.Name,
		Email:        payload.Email,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleGetInvites(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	invites, err := h.storage.Query.ListOrganizationInvites(r.Context(), pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, invites)
}
