package auth

import (
	"context"
	"errors"
	"time"

	"github.com/mixdone/uptime-monitoring/internal/models/dto"
	"github.com/mixdone/uptime-monitoring/internal/models/errs"
	"github.com/mixdone/uptime-monitoring/internal/services/constants"
	"github.com/mixdone/uptime-monitoring/internal/services/interfaces"
	"github.com/mixdone/uptime-monitoring/pkg/logger"
)

type authService struct {
	logger  logger.Logger
	user    interfaces.UserService
	session interfaces.SessionService
	token   interfaces.TokenService
}

func NewAuthService(user interfaces.UserService, session interfaces.SessionService, token interfaces.TokenService, log logger.Logger) interfaces.AuthenticationService {
	return &authService{
		logger:  log,
		user:    user,
		session: session,
		token:   token,
	}
}

func (a *authService) Register(ctx context.Context, username, password, fingerprint string) (*dto.AuthResult, error) {
	id, err := a.user.RegisterUser(ctx, username, password)
	if err != nil {
		return nil, err
	}

	return a.createAuthResult(ctx, id, fingerprint)

}

func (a *authService) Login(ctx context.Context, username, password, fingerprint string) (*dto.AuthResult, error) {
	user, err := a.user.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if !a.user.VerifyPassword(user.PasswordHash, password) {
		return nil, errors.New("wrong password")
	}

	return a.createAuthResult(ctx, user.ID, fingerprint)
}

func (a *authService) Logout(ctx context.Context, userID int, refreshToken, fingerprint string) error {
	session, err := a.session.GetSession(ctx, userID, refreshToken, fingerprint)
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

func (a *authService) RefreshTokens(ctx context.Context, userID int, refreshToken, fingerprint string) (*dto.AuthResult, error) {

	err := a.Logout(ctx, userID, refreshToken, fingerprint)
	if err != nil {
		return nil, err
	}

	return a.createAuthResult(ctx, userID, fingerprint)
}

func (a *authService) createAuthResult(ctx context.Context, userID int, fingerprint string) (*dto.AuthResult, error) {
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
