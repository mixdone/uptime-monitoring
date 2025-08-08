package dto

import "encoding/json"

type MonitorRequest struct {
	Name             string          `json:"name" binding:"required"`
	Type             string          `json:"type" binding:"required"`
	Target           string          `json:"target" binding:"required,url"`
	Timeout          int             `json:"timeout" binding:"required,gte=0"`
	Interval         int             `json:"interval" binding:"required,gte=0"`
	IsActive         bool            `json:"is_active"`
	RequestSpec      json.RawMessage `json:"request_spec" binding:"required"`
	ExpectedResponse json.RawMessage `json:"expected_response"`
}

type MonitorResponse struct {
	ID int64 `json:"id"`
}
