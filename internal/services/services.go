package services

import (
	"context"

	"github.com/mixdone/uptime-monitoring/internal/models"
	"github.com/mixdone/uptime-monitoring/internal/repository"
)

type UserService interface {
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByID(ctx context.Context, userID int) (*models.User, error)
	VerifyPassword(hashFromDB, inputPassword string) bool
	RegisterUser(ctx context.Context, username, password string) (int, error)
	DeleteUser(ctx context.Context, userID int) error
}

type AuthorizationService interface{}

type Services struct {
	User UserService
	Auth AuthorizationService
}

func NewServices(repositories *repository.Repository) *Services {
	return nil
}
