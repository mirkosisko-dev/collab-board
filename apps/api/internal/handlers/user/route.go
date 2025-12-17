package user

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mirkosisko-dev/api/config"
	pool "github.com/mirkosisko-dev/api/db"
	"github.com/mirkosisko-dev/api/db/sqlc"
	"github.com/mirkosisko-dev/api/internal/handlers/auth"
	"github.com/mirkosisko-dev/api/internal/services/session"
	"github.com/mirkosisko-dev/api/internal/types"
	"github.com/mirkosisko-dev/api/utils"
)

type Handler struct {
	storage        *pool.Database
	config         *config.Config
	sessionService *session.Service
}

func NewHandler(storage *pool.Database, config *config.Config, sessionService *session.Service) *Handler {
	return &Handler{
		storage:        storage,
		config:         config,
		sessionService: sessionService,
	}
}

func (h *Handler) RegisterPublicRoutes(r *mux.Router) {
	r.HandleFunc("/auth/register", h.handleRegister).Methods(http.MethodPost)
	r.HandleFunc("/auth/login", h.handleLogin).Methods(http.MethodPost)
	r.HandleFunc("/auth/logout", h.handleLogout).Methods(http.MethodPost)
	r.HandleFunc("/auth/refresh", h.renewAccessToken).Methods(http.MethodPost)
}

func (h *Handler) RegisterProtectedRoutes(r *mux.Router) {
	r.HandleFunc("/me", h.handleGetMe).Methods(http.MethodGet)
	// r.HandleFunc("/users/me/password", h.handleUpdatePassword).Methods(http.MethodPatch)
	r.HandleFunc("/invites", h.handleGetInvites).Methods(http.MethodGet)
	r.HandleFunc("/auth/revoke", h.renewAccessToken).Methods(http.MethodPost)
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

	if !user.ID.Valid {
		utils.WriteError(w, http.StatusInternalServerError, errors.New("user id invalid"))
		return
	}

	userUUID := uuid.UUID(user.ID.Bytes)

	session, accessToken, refreshToken, atExp, rtExp, err := h.sessionService.CreateSession(r.Context(), userUUID, user.Email)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, types.LoginUserRes{
		SessionID:             session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  atExp.String(),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: rtExp.String(),
		User: types.UserRes{
			Name:  user.Name,
			Email: user.Email,
		},
	})
}

func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(auth.AuthKey{}).(*auth.TokenClaims)

	sessionUUID, err := uuid.Parse(claims.RegisteredClaims.ID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = h.sessionService.DeleteSession(r.Context(), sessionUUID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
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

func (h *Handler) handleGetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	u, err := h.storage.Query.GetUser(r.Context(), pgtype.UUID{
		Bytes: userID,
		Valid: true,
	})
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	resp := types.UserRes{
		Name:  u.Name,
		Email: u.Email,
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) renewAccessToken(w http.ResponseWriter, r *http.Request) {
	var payload types.RenewAccessTokenPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	accessToken, atExp, err := h.sessionService.RenewAccessToken(r.Context(), payload.RefreshToken)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, types.RenewAccessTokenRes{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: atExp.String(),
	})
}
