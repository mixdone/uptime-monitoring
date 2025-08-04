package services

import (
	"github.com/mixdone/uptime-monitoring/internal/config"
	"github.com/mixdone/uptime-monitoring/internal/repository"
	"github.com/mixdone/uptime-monitoring/internal/services/auth"
	"github.com/mixdone/uptime-monitoring/internal/services/constants"
	"github.com/mixdone/uptime-monitoring/internal/services/session"
	"github.com/mixdone/uptime-monitoring/internal/services/token"
	"github.com/mixdone/uptime-monitoring/internal/services/user"
	"github.com/mixdone/uptime-monitoring/pkg/logger"
)

type Services struct {
	User    user.UserService
	Token   token.TokenService
	Session session.SessionService
	Auth    auth.AuthenticationService
}

func NewServices(repositories *repository.Repository, cfg config.Config, log logger.Logger) *Services {
	user := user.NewUserService(repositories.Users, log)
	token := token.NewTokenService(cfg.Jwt.AccessSecret, cfg.Jwt.RefreshSecret, constants.AccessTokenTTL, constants.RefreshTokenTTL)
	session := session.NewSessionService(repositories.Sessions, log)
	auth := auth.NewAuthService(user, session, token, log)

	return &Services{
		User:    user,
		Token:   token,
		Session: session,
		Auth:    auth,
	}
}
