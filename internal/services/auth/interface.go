package auth

import (
	"context"

	"github.com/mixdone/uptime-monitoring/internal/models/dto"
)

type AuthenticationService interface {
	Register(ctx context.Context, userDTO dto.RegisterRequest) (*dto.AuthResult, error)
	Login(ctx context.Context, userDTO dto.LoginRequest) (*dto.AuthResult, error)
	Logout(ctx context.Context, userID int64, userDTO dto.LogoutRequest) error
	RefreshTokens(ctx context.Context, userID int64, userDTO dto.RefreshRequest) (*dto.AuthResult, error)
}
