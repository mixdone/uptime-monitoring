package models

import (
	"encoding/json"
	"time"
)

type Monitor struct {
	ID     int64 `json:"id" db:"id"`
	UserID int64 `json:"user_id" db:"user_id"`

	Name          string     `json:"name" db:"name"`
	Type          string     `json:"type" db:"type"`
	Target        string     `json:"target" db:"target"`
	Timeout       int        `json:"timeout" db:"timeout"`
	Interval      int        `json:"interval" db:"interval"`
	IsActive      bool       `json:"is_active" db:"is_active"`
	LastCheckedAt *time.Time `json:"last_checked_at,omitempty" db:"last_checked_at"`

	RequestSpec      json.RawMessage `json:"request" db:"request"`
	ExpectedResponse json.RawMessage `json:"expected_response" db:"expected_response"`
}
