package models

import (
	"encoding/json"
	"fmt"
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

	RequestSpec      json.RawMessage `json:"request" db:"request" swaggerignore:"true"`
	ExpectedResponse json.RawMessage `json:"expected_response" db:"expected_response" swaggerignore:"true"`
}

type HTTPMonitorTask struct {
	ID      int64  `json:"id"`
	URL     string `json:"target"`
	Method  string `json:"method"`
	Timeout int    `json:"timeout"`

	Headers     map[string]string `json:"headers,omitempty"`
	RequestBody string            `json:"body,omitempty"`

	ExpectStatus       int      `json:"expect_status"`
	ExpectBodyContains []string `json:"expect_body_contains,omitempty"`
}

type CheckResult struct {
	TaskID    int64     `json:"task_id"`
	Errors    []string  `json:"errors"`
	Duration  int64     `json:"duration_ms"`
	CheckedAt time.Time `json:"checked_at"`
}

func (m *Monitor) ToHTTPMonitorTask() (*HTTPMonitorTask, error) {
	var spec struct {
		Method  string            `json:"method"`
		Headers map[string]string `json:"headers,omitempty"`
		Body    string            `json:"body,omitempty"`
	}
	if err := json.Unmarshal(m.RequestSpec, &spec); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request spec: %w", err)
	}

	var expected struct {
		Status       int      `json:"expect_status"`
		BodyContains []string `json:"expect_body_contains,omitempty"`
	}
	if len(m.ExpectedResponse) > 0 {
		if err := json.Unmarshal(m.ExpectedResponse, &expected); err != nil {
			return nil, fmt.Errorf("failed to unmarshal expected response: %w", err)
		}
	}

	return &HTTPMonitorTask{
		ID:      m.ID,
		URL:     m.Target,
		Method:  spec.Method,
		Timeout: m.Timeout,

		Headers:            spec.Headers,
		RequestBody:        spec.Body,
		ExpectStatus:       expected.Status,
		ExpectBodyContains: expected.BodyContains,
	}, nil
}
