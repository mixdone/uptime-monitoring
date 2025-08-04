package session

import (
	"context"
	"time"

	"github.com/mixdone/uptime-monitoring/internal/models"
)

type SessionService interface {
	CreateSession(ctx context.Context, userID int64,
		refreshToken, fingerprint string, expiresAt time.Time) (int64, error)
	GetSession(ctx context.Context, userID int64,
		refreshToken, fingerprint string) (*models.Session, error)
	DeleteSession(ctx context.Context, sessionID int64) error
	DeleteAllUserSessions(ctx context.Context, userID int64) error
}
