package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mixdone/uptime-monitoring/internal/models"
	"github.com/mixdone/uptime-monitoring/internal/models/errs"
)

type sessionRepository struct {
	db *pgxpool.Pool
}

func NewSessionRepo(pool *pgxpool.Pool) SessionRepository {
	return &sessionRepository{db: pool}
}

func (s *sessionRepository) CreateSession(ctx context.Context, session models.Session) (int64, error) {
	var id int64
	query := `
		INSERT INTO sessions (user_id, refresh_token, expires_at, fingerprint)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	err := s.db.QueryRow(ctx, query, session.UserID,
		session.RefreshToken, session.ExpiresAt, session.Fingerprint).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *sessionRepository) GetSession(ctx context.Context, userID int, refreshToken, fingerprint string) (*models.Session, error) {
	var session models.Session

	query := `
		SELECT id, user_id, refresh_token, expires_at, fingerprint
		FROM sessions
		WHERE user_id = $1 AND refresh_token = $2 AND fingerprint = $3`

	err := s.db.QueryRow(ctx, query, userID, refreshToken, fingerprint).Scan(
		&session.ID,
		&session.UserID,
		&session.RefreshToken,
		&session.ExpiresAt,
		&session.Fingerprint,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errs.ErrSessionNotFound
	}

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *sessionRepository) GetAllUserSessions(ctx context.Context, userID int) ([]models.Session, error) {
	var sessions []models.Session

	query := `
		SELECT id, user_id, refresh_token, fingerprint
		FROM sessions
		WHERE user_id = $1
	`

	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var session models.Session
		if err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.RefreshToken,
			&session.ExpiresAt,
			&session.Fingerprint,
		); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		sessions = append(sessions, session)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iteration error: %w", err)
	}

	return sessions, nil
}

func (s *sessionRepository) DeleteSession(ctx context.Context, sessionID int64) error {
	query := `
		DELETE FROM sessions 
		WHERE id = $1
	`
	_, err := s.db.Exec(ctx, query, sessionID)

	return err
}

func (s *sessionRepository) DeleteAllSessions(ctx context.Context, userID int) error {
	query := `
		DELETE FROM sessions 
		WHERE user_id = $1
	`
	_, err := s.db.Exec(ctx, query, userID)

	return err
}
