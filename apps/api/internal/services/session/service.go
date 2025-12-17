package session

import (
	"context"
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

func (s *Service) CreateSession(ctx context.Context, userID uuid.UUID) (*sqlc.Session, string, string, time.Time, time.Time, error) {
	// 1. Create Access Token
	accessToken, err := auth.CreateJWT(userID, s.config.JWTSecret, s.config.JWTExpirationInSeconds)
	if err != nil {
		return nil, "", "", time.Time{}, time.Time{}, err
	}

	// 2. Create Refresh Token
	refreshToken, err := auth.CreateRefreshToken(userID, s.config.RefreshTokenSecret, s.config.RefreshTokenExpirationInSeconds)
	if err != nil {
		return nil, "", "", time.Time{}, time.Time{}, err
	}

	// Calc expirations
	rtExp := time.Now().Add(time.Second * time.Duration(s.config.RefreshTokenExpirationInSeconds))
	atExp := time.Now().Add(time.Second * time.Duration(s.config.JWTExpirationInSeconds))

	// 3. Store in DB
	session, err := s.storage.Query.CreateSesion(ctx, sqlc.CreateSesionParams{
		UserID:       pgtype.UUID{Bytes: userID, Valid: true},
		RefreshToken: refreshToken,
		IsRevoked:    false,
		ExpiresAt:    pgtype.Timestamp{Time: rtExp, Valid: true},
	})
	if err != nil {
		return nil, "", "", time.Time{}, time.Time{}, err
	}

	return &session, accessToken, refreshToken, atExp, rtExp, nil
}
