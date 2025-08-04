package auth

import (
	"context"
	"errors"
	"time"

	"github.com/mixdone/uptime-monitoring/internal/models/dto"
	"github.com/mixdone/uptime-monitoring/internal/models/errs"
	"github.com/mixdone/uptime-monitoring/internal/services/constants"
	"github.com/mixdone/uptime-monitoring/internal/services/session"
	"github.com/mixdone/uptime-monitoring/internal/services/token"
	"github.com/mixdone/uptime-monitoring/internal/services/user"
	"github.com/mixdone/uptime-monitoring/pkg/logger"
)

type authService struct {
	logger  logger.Logger
	user    user.UserService
	session session.SessionService
	token   token.TokenService
}

func NewAuthService(user user.UserService, session session.SessionService, token token.TokenService, log logger.Logger) AuthenticationService {
	return &authService{
		logger:  log,
		user:    user,
		session: session,
		token:   token,
	}
}

func (a *authService) Register(ctx context.Context, userDTO dto.RegisterRequest) (*dto.AuthResult, error) {
	id, err := a.user.RegisterUser(ctx, userDTO)
	if err != nil {
		return nil, err
	}

	return a.createAuthResult(ctx, id, userDTO.Fingerprint)

}

func (a *authService) Login(ctx context.Context, userDTO dto.LoginRequest) (*dto.AuthResult, error) {
	user, err := a.user.GetByUsername(ctx, userDTO.Username)
	if err != nil {
		return nil, err
	}

	if !a.user.VerifyPassword(user.PasswordHash, userDTO.Password) {
		return nil, errors.New("wrong password")
	}

	return a.createAuthResult(ctx, user.ID, userDTO.Fingerprint)
}

func (a *authService) Logout(ctx context.Context, userID int64, userDTO dto.LogoutRequest) error {
	session, err := a.session.GetSession(ctx, userID, userDTO.RefreshToken, userDTO.Fingerprint)
	if errors.Is(err, errs.ErrSessionNotFound) {
		return nil
	} else if err != nil {
		return err
	}

	if err = a.session.DeleteSession(ctx, session.ID); err != nil {
		return err
	}

	return nil
}

func (a *authService) RefreshTokens(ctx context.Context, userID int64, userDTO dto.RefreshRequest) (*dto.AuthResult, error) {

	err := a.Logout(ctx, userID, dto.LogoutRequest(userDTO))
	if err != nil {
		return nil, err
	}

	return a.createAuthResult(ctx, userID, userDTO.Fingerprint)
}

func (a *authService) createAuthResult(ctx context.Context, userID int64, fingerprint string) (*dto.AuthResult, error) {
	accessToken, refreshToken, err := a.token.Generate(userID)
	if err != nil {
		return nil, err
	}

	expireAt := time.Now().Add(constants.RefreshTokenTTL)
	_, err = a.session.CreateSession(ctx, userID, refreshToken, fingerprint, expireAt)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
