package monitors

import (
	"context"
	"time"

	"github.com/mixdone/uptime-monitoring/internal/models"
)

type MonitorService interface {
	CreateMonitor(ctx context.Context, monitor models.Monitor) (int64, error)
	GetMonitor(ctx context.Context, id int64) (*models.Monitor, error)
	GetAllUserMonitors(ctx context.Context, userID int64) ([]models.Monitor, error)
	GetAllActiveMonitors(ctx context.Context) ([]models.Monitor, error)
	UpdateMonitor(ctx context.Context, monitor models.Monitor) error
	UpdateLastCheckedAt(ctx context.Context, id int64, checkedAt time.Time) error
	DeleteMonitor(ctx context.Context, id int64) error
}
