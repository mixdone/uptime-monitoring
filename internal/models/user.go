package models

import "time"

type User struct {
	ID           int    `json:"id" db:"id"`
	Username     string `json:"username" db:"username"`
	Email        string `json:"email,omitempty" db:"email"`
	TelegramID   int64  `json:"telegram_id" db:"telegram_id"`
	PasswordHash string `json:"-" db:"password_hash"`
}

type Session struct {
	ID           int       `json:"id" db:"id"`
	UserID       int       `json:"user_id" db:"user_id"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token"`
	ExpiresAt    time.Time `json:"-" db:"expires_at"`
	Fingerprint  string    `json:"fingerprint" db:"fingerprint"`
}
