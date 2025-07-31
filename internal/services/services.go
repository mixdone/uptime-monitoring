package services

import (
	"context"
	"time"

	"github.com/mixdone/uptime-monitoring/internal/models"
	"github.com/mixdone/uptime-monitoring/internal/models/dto"
	"github.com/mixdone/uptime-monitoring/internal/repository"
)

type UserService interface {
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByID(ctx context.Context, userID int) (*models.User, error)
	VerifyPassword(hashFromDB, inputPassword string) bool
	RegisterUser(ctx context.Context, username, password string) (int, error)
	DeleteUser(ctx context.Context, userID int) error
}

type SessionService interface {
	CreateSession(ctx context.Context, userID int,
		refreshToken, fingerprint string, expiresAt time.Time) (int64, error)
	GetSession(ctx context.Context, userID int,
		refreshToken, fingerprint string) (*models.Session, error)
	DeleteSession(ctx context.Context, sessionID int64) error
	DeleteAllUserSessions(ctx context.Context, userID int) error
}

type TokenService interface {
	Generate(userID int) (accessToken, refreshToken string, err error)
	ValidateAccess(tokenStr string) (userID int, err error)
	ValidateRefresh(tokenStr string) (userID int, err error)
}

type AuthenticationService interface {
	Register(ctx context.Context, username, password, fingerprint string) (*dto.AuthResult, error)
	Login(ctx context.Context, username, password, fingerprint string) (*dto.AuthResult, error)
	Logout(ctx context.Context, userID int, password, fingerprint string) error
	RefreshTokens(ctx context.Context, userID int, refreshToken, fingerprint string) (*dto.AuthResult, error)
}

type Services struct {
	User UserService
	Auth AuthenticationService
}

func NewServices(repositories *repository.Repository) *Services {
	return nil
}
