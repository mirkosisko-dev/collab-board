package session

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mirkosisko-dev/api/config"
	pool "github.com/mirkosisko-dev/api/db"
	"github.com/mirkosisko-dev/api/db/sqlc"
	"github.com/mirkosisko-dev/api/internal/handlers/auth"
)

type Service struct {
	storage *pool.Database
	config  *config.Config
}

func NewService(storage *pool.Database, cfg *config.Config) *Service {
	return &Service{
		storage: storage,
		config:  cfg,
	}
}

func (s *Service) CreateSession(ctx context.Context, userID uuid.UUID, email string) (*sqlc.Session, string, string, time.Time, time.Time, error) {
	accessToken, atClaims, err := auth.CreateAccessToken(userID, email, s.config.AccessTokenSecret, s.config.AccessTokenExpirationInSeconds)
	if err != nil {
		return nil, "", "", time.Time{}, time.Time{}, err
	}

	refreshToken, rtClaims, err := auth.CreateRefreshToken(userID, email, s.config.RefreshTokenSecret, s.config.RefreshTokenExpirationInSeconds)
	if err != nil {
		return nil, "", "", time.Time{}, time.Time{}, err
	}

	session, err := s.storage.Query.CreateSesion(ctx, sqlc.CreateSesionParams{
		ID:           pgtype.UUID{Bytes: rtClaims.UserID, Valid: true},
		RefreshToken: refreshToken,
		IsRevoked:    false,
		ExpiresAt:    pgtype.Timestamp{Time: rtClaims.ExpiresAt.Time, Valid: true},
	})
	if err != nil {
		return nil, "", "", time.Time{}, time.Time{}, err
	}

	return &session, accessToken, refreshToken, atClaims.ExpiresAt.Time, rtClaims.ExpiresAt.Time, nil
}

func (s *Service) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	err := s.storage.Query.DeleteSession(ctx, pgtype.UUID{Bytes: sessionID, Valid: true})
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	err := s.storage.Query.RevokeSession(ctx, pgtype.UUID{Bytes: sessionID, Valid: true})
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) RenewAccessToken(ctx context.Context, refreshToken string) (string, time.Time, error) {
	rtClaims, err := auth.ValidateToken(refreshToken, s.config.RefreshTokenSecret)
	if err != nil {
		return "", time.Time{}, errors.New("invalid refresh token")
	}

	userID, err := uuid.Parse(rtClaims.UserID.String())
	if err != nil {
		return "", time.Time{}, errors.New("invalid user id")
	}

	sessionUUID, err := uuid.Parse(rtClaims.RegisteredClaims.ID)
	if err != nil {
		return "", time.Time{}, errors.New("error parsing string")
	}

	session, err := s.storage.Query.GetSession(ctx, pgtype.UUID{Bytes: sessionUUID, Valid: true})
	if err != nil {
		return "", time.Time{}, errors.New("session not found")
	}

	if session.IsRevoked {
		return "", time.Time{}, errors.New("session revoked")
	}

	if session.ExpiresAt.Time.Before(time.Now()) {
		return "", time.Time{}, errors.New("session expired")
	}

	if session.ID.Bytes != userID {
		return "", time.Time{}, errors.New("session user mismatch")
	}

	accessToken, atClaims, err := auth.CreateAccessToken(userID, rtClaims.Email, s.config.AccessTokenSecret, s.config.AccessTokenExpirationInSeconds)
	if err != nil {
		return "", time.Time{}, err
	}

	return accessToken, atClaims.ExpiresAt.Time, nil
}
