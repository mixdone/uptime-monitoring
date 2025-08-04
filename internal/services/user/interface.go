package user

import (
	"context"

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
