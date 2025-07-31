package user

import (
	"context"
	"errors"

	"github.com/mixdone/uptime-monitoring/internal/models"
	"github.com/mixdone/uptime-monitoring/internal/models/errs"
	"github.com/mixdone/uptime-monitoring/internal/repository"
	"github.com/mixdone/uptime-monitoring/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo   repository.UserRepository
	logger logger.Logger
}

func NewUserService(repo repository.UserRepository, log logger.Logger) *userService {
	return &userService{
		repo:   repo,
		logger: log.WithField("component", "userService"),
	}
}

func (s *userService) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	s.logger.Debugf("Fetching user by username: %s", username)

	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		s.logger.WithField("username", username).
			WithError(err).
			Error("Failed to get user by username")
		return nil, err
	}

	return user, nil
}

func (s *userService) GetByID(ctx context.Context, userID int) (*models.User, error) {
	s.logger.Debugf("Fetching user by ID: %d", userID)

	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		s.logger.WithField("userID", userID).
			WithError(err).
			Error("Failed to get user by ID")
		return nil, err
	}

	return user, nil
}

func (s *userService) RegisterUser(ctx context.Context, username, password string) (int, error) {
	s.logger.Infof("Attempting to register user: %s", username)

	_, err := s.repo.GetUserByUsername(ctx, username)
	if err == nil {
		s.logger.Warnf("Username already taken: %s", username)
		return 0, errs.ErrUsernameTaken
	}

	if !errors.Is(err, errs.ErrUserNotFound) {
		s.logger.WithError(err).Error("Unexpected error while checking username")
		return 0, errs.ErrInternal
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.WithField("username", username).
			WithError(err).
			Error("Failed to hash password")
		return 0, errs.ErrHashingFailed
	}

	user := models.User{
		Username:     username,
		PasswordHash: string(hash),
	}

	id, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		s.logger.WithField("username", username).
			WithError(err).
			Error("Failed to create user in DB")
		return 0, err
	}

	s.logger.Infof("User registered successfully: %s (id=%d)", username, id)
	return id, nil
}

func (s *userService) VerifyPassword(hashFromDB, inputPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashFromDB), []byte(inputPassword))
	if err != nil {
		s.logger.Debug("Password verification failed")
		return false
	}
	return true
}

func (s *userService) DeleteUser(ctx context.Context, userID int) error {
	s.logger.Infof("Deleting user: id=%d", userID)

	if err := s.repo.DeleteUser(ctx, userID); err != nil {
		s.logger.WithField("userID", userID).
			WithError(err).
			Error("Failed to delete user")
		return err
	}

	s.logger.Infof("User deleted successfully: id=%d", userID)
	return nil
}
