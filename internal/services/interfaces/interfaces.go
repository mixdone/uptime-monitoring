package interfaces

import (
	"context"
	"time"

	"github.com/mixdone/uptime-monitoring/internal/models"
	"github.com/mixdone/uptime-monitoring/internal/models/dto"
)

type UserService interface {
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByID(ctx context.Context, userID int64) (*models.User, error)
	VerifyPassword(hashFromDB, inputPassword string) bool
	RegisterUser(ctx context.Context, userDTO dto.RegisterRequest) (int64, error)
	DeleteUser(ctx context.Context, userID int64) error
}

type SessionService interface {
	CreateSession(ctx context.Context, userID int64,
		refreshToken, fingerprint string, expiresAt time.Time) (int64, error)
	GetSession(ctx context.Context, userID int64,
		refreshToken, fingerprint string) (*models.Session, error)
	DeleteSession(ctx context.Context, sessionID int64) error
	DeleteAllUserSessions(ctx context.Context, userID int64) error
}

type TokenService interface {
	Generate(userID int64) (accessToken, refreshToken string, err error)
	ValidateAccess(tokenStr string) (userID int64, err error)
	ValidateRefresh(tokenStr string) (userID int64, err error)
}

type AuthenticationService interface {
	Register(ctx context.Context, userDTO dto.RegisterRequest) (*dto.AuthResult, error)
	Login(ctx context.Context, userDTO dto.LoginRequest) (*dto.AuthResult, error)
	Logout(ctx context.Context, userID int64, userDTO dto.LogoutRequest) error
	RefreshTokens(ctx context.Context, userID int64, userDTO dto.RefreshRequest) (*dto.AuthResult, error)
}
