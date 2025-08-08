package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mixdone/uptime-monitoring/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) (int64, error)
	GetUser(ctx context.Context, userId int64) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	//UpdateUser(ctx context.Context, user models.User) error
	DeleteUser(ctx context.Context, userId int64) error
}

type SessionRepository interface {
	CreateSession(ctx context.Context, session models.Session) (int64, error)
	GetSession(ctx context.Context, userID int64, refreshToken, fingerprint string) (*models.Session, error)
	GetAllUserSessions(ctx context.Context, userID int64) ([]models.Session, error)
	DeleteSession(ctx context.Context, sessionID int64) error
	DeleteAllSessions(ctx context.Context, userID int64) error
}

type MonitorsRepository interface {
	CreateMonitor(ctx context.Context, monitor models.Monitor) (int64, error)
	GetMonitor(ctx context.Context, id int64) (*models.Monitor, error)
	GetAllUserMonitors(ctx context.Context, userID int64) ([]models.Monitor, error)
	GetAllActiveMonitors(ctx context.Context) ([]models.Monitor, error)
	UpdateMonitor(ctx context.Context, monitor models.Monitor) error
	UpdateLastCheckedAt(ctx context.Context, id int64, checkedAt time.Time) error
	DeleteMonitor(ctx context.Context, id int64) error
}

type Repository struct {
	Users    UserRepository
	Sessions SessionRepository
	Monitors MonitorsRepository
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Users:    NewUserRepo(db),
		Sessions: NewSessionRepo(db),
		Monitors: NewMonitorRepo(db),
	}
}
