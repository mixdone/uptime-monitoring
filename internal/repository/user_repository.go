package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mixdone/uptime-monitoring/internal/models"
	"github.com/mixdone/uptime-monitoring/internal/models/errs"
)

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) UserRepository {
	return &userRepo{db: pool}
}

func (u *userRepo) CreateUser(ctx context.Context, user models.User) (int, error) {
	var id int
	query := `
		INSERT INTO users (username, email, telegram_id, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	err := u.db.QueryRow(ctx, query,
		user.Username, user.Email,
		user.TelegramID, user.PasswordHash).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (u *userRepo) GetUser(ctx context.Context, userId int) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, username, email, telegram_id, password_hash 
		FROM users
		WHERE id = $1`

	err := u.db.QueryRow(ctx, query, userId).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.TelegramID,
		&user.PasswordHash,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (u *userRepo) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	query := `
			SELECT id, username, email, telegram_id, password_hash
			FROM users
			WHERE username = $1`

	err := u.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.TelegramID,
		&user.PasswordHash,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (u *userRepo) DeleteUser(ctx context.Context, userId int) error {
	query := `
		DELETE FROM users 
		WHERE id = $1 
	`

	cmdTag, err := u.db.Exec(ctx, query, userId)

	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errs.ErrUserNotFound
	}

	return nil
}
