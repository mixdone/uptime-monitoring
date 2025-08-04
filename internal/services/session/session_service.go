package session

import (
	"context"
	"errors"
	"time"

	"github.com/mixdone/uptime-monitoring/internal/models"
	"github.com/mixdone/uptime-monitoring/internal/models/errs"
	"github.com/mixdone/uptime-monitoring/internal/repository"
	"github.com/mixdone/uptime-monitoring/internal/services/interfaces"
	"github.com/mixdone/uptime-monitoring/pkg/logger"
)

type sessionService struct {
	repo   repository.SessionRepository
	logger logger.Logger
}

func NewSessionService(repo repository.SessionRepository, log logger.Logger) interfaces.SessionService {
	return &sessionService{
		repo:   repo,
		logger: log.WithField("component", "sessionService"),
	}
}

func (s *sessionService) CreateSession(ctx context.Context, userID int64,
	refreshToken, fingerprint string,
	expiresAt time.Time) (int64, error) {
	session := models.Session{
		UserID:       userID,
		RefreshToken: refreshToken,
		Fingerprint:  fingerprint,
		ExpiresAt:    expiresAt,
	}

	id, err := s.repo.CreateSession(ctx, session)

	if err != nil {
		s.logger.WithField("user_id", userID).
			WithError(err).
			Error("Failed to create session")
		return 0, err
	}

	s.logger.WithFields(map[string]any{
		"user_id":     userID,
		"session":     id,
		"fingerprint": fingerprint,
	}).Info("Session created successfully")

	return id, nil
}

func (s *sessionService) GetSession(ctx context.Context, userID int64,
	refreshToken, fingerprint string) (*models.Session, error) {

	session, err := s.repo.GetSession(ctx, userID, refreshToken, fingerprint)

	if errors.Is(err, errs.ErrSessionNotFound) {
		s.logger.WithFields(map[string]any{
			"user_id":     userID,
			"fingerprint": fingerprint,
		}).WithError(err).Info("Session not found")
		return nil, err
	} else if err != nil {
		s.logger.WithFields(map[string]any{
			"user_id":     userID,
			"fingerprint": fingerprint,
		}).WithError(err).Error("Repository error")
		return nil, err
	}

	s.logger.WithFields(map[string]any{
		"user_id":     userID,
		"session":     session.ID,
		"fingerprint": fingerprint,
	}).Info("Session found successfully")

	return session, nil
}

func (s *sessionService) DeleteSession(ctx context.Context, sessionID int64) error {
	err := s.repo.DeleteSession(ctx, sessionID)
	if err != nil {
		s.logger.WithField("session_id", sessionID).
			WithError(err).
			Error("Failed to delete session")
		return err
	}
	s.logger.WithField("session_id", sessionID).
		Info("Session deleted successfully")
	return nil
}

func (s *sessionService) DeleteAllUserSessions(ctx context.Context, userID int64) error {
	err := s.repo.DeleteAllSessions(ctx, userID)
	if err != nil {
		s.logger.WithField("user_id", userID).
			WithError(err).
			Error("Failed to delete all sessions for user")
		return err
	}
	s.logger.WithField("user_id", userID).
		Info("All sessions deleted for user")
	return nil
}
