package monitors

import (
	"context"
	"time"

	"github.com/mixdone/uptime-monitoring/internal/models"
	"github.com/mixdone/uptime-monitoring/internal/repository"
	"github.com/mixdone/uptime-monitoring/pkg/logger"
)

type monitorService struct {
	repo   repository.MonitorsRepository
	logger logger.Logger
}

func NewMonitorService(repo repository.MonitorsRepository, log logger.Logger) MonitorService {
	return &monitorService{
		repo:   repo,
		logger: log.WithField("component", "monitorService"),
	}
}

func (s *monitorService) CreateMonitor(ctx context.Context, monitor models.Monitor) (int64, error) {
	s.logger.Infof("Creating monitor for user_id=%d name=%s", monitor.UserID, monitor.Name)

	id, err := s.repo.CreateMonitor(ctx, monitor)
	if err != nil {
		s.logger.WithFields(map[string]any{
			"userID":      monitor.UserID,
			"monitorName": monitor.Name,
		}).WithError(err).Error("Failed to create monitor")
		return 0, err
	}

	s.logger.Infof("Monitor created successfully with id=%d", id)
	return id, nil
}

func (s *monitorService) GetMonitor(ctx context.Context, id int64) (*models.Monitor, error) {
	s.logger.Debugf("Fetching monitor with id=%d", id)

	monitor, err := s.repo.GetMonitor(ctx, id)
	if err != nil {
		s.logger.WithFields(map[string]any{
			"monitorID": id,
		}).WithError(err).Error("Failed to fetch monitor")
		return nil, err
	}

	return monitor, nil
}

func (s *monitorService) GetAllUserMonitors(ctx context.Context, userID int64) ([]models.Monitor, error) {
	s.logger.Debugf("Fetching all monitors for user_id=%d", userID)

	monitors, err := s.repo.GetAllUserMonitors(ctx, userID)
	if err != nil {
		s.logger.WithFields(map[string]any{
			"userID": userID,
		}).WithError(err).Error("Failed to fetch user monitors")
		return nil, err
	}

	return monitors, nil
}

func (s *monitorService) GetAllActiveMonitors(ctx context.Context) ([]models.Monitor, error) {
	s.logger.Debug("Fetching all active monitors")

	monitors, err := s.repo.GetAllActiveMonitors(ctx)
	if err != nil {
		s.logger.WithError(err).Error("Failed to fetch active monitors")
		return nil, err
	}

	return monitors, nil
}

func (s *monitorService) UpdateMonitor(ctx context.Context, monitor models.Monitor) error {
	s.logger.Infof("Updating monitor id=%d", monitor.ID)

	if err := s.repo.UpdateMonitor(ctx, monitor); err != nil {
		s.logger.WithFields(map[string]any{
			"monitorID": monitor.ID,
		}).WithError(err).Error("Failed to update monitor")
		return err
	}

	s.logger.Infof("Monitor updated successfully id=%d", monitor.ID)
	return nil
}

func (s *monitorService) UpdateLastCheckedAt(ctx context.Context, id int64, checkedAt time.Time) error {
	s.logger.Debugf("Updating last_checked_at for monitor id=%d", id)

	if err := s.repo.UpdateLastCheckedAt(ctx, id, checkedAt); err != nil {
		s.logger.WithFields(map[string]any{
			"monitorID": id,
		}).WithError(err).Error("Failed to update last_checked_at")
		return err
	}

	s.logger.Debugf("Updated last_checked_at successfully for monitor id=%d", id)
	return nil
}

func (s *monitorService) DeleteMonitor(ctx context.Context, id int64) error {
	s.logger.Infof("Deleting monitor id=%d", id)

	if err := s.repo.DeleteMonitor(ctx, id); err != nil {
		s.logger.WithFields(map[string]any{
			"monitorID": id,
		}).WithError(err).Error("Failed to delete monitor")
		return err
	}

	s.logger.Infof("Monitor deleted successfully id=%d", id)
	return nil
}
