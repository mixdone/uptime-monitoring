package auth

import (
	"context"
	"time"

	"github.com/mixdone/uptime-monitoring/internal/models"
)

type TokenService interface {
	Generate(userID int) (accessToken, refreshToken string, err error)
	ValidateAccess(tokenStr string) (userID int, err error)
	ValidateRefresh(tokenStr string) (userID int, err error)
}

type SessionService interface {
	CreateSession(ctx context.Context, userID int,
		refreshToken, fingerprint string, expiresAt time.Time) (int64, error)
	GetSession(ctx context.Context, userID int,
		refreshToken, fingerprint string) (*models.Session, error)
	DeleteSession(ctx context.Context, sessionID int64) error
	DeleteAllUserSessions(ctx context.Context, userID int) error
}
