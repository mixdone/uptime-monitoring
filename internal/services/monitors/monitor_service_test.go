package monitors_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"

	"github.com/mixdone/uptime-monitoring/internal/mocks"
	"github.com/mixdone/uptime-monitoring/internal/models"
	"github.com/mixdone/uptime-monitoring/internal/services/monitors"
)

const expectedID = int64(123456789)

func setup(t *testing.T) (context.Context, *gomock.Controller, *mocks.MockMonitorsRepository, *mocks.MockLogger, monitors.MonitorService) {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockMonitorsRepository(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)

	mockLogger.EXPECT().WithField(gomock.Any(), gomock.Any()).Return(mockLogger).AnyTimes()
	mockLogger.EXPECT().WithFields(gomock.Any()).Return(mockLogger).AnyTimes()
	mockLogger.EXPECT().WithError(gomock.Any()).Return(mockLogger).AnyTimes()

	mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Debugf(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Info(gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Debug(gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any()).AnyTimes()

	svc := monitors.NewMonitorService(mockRepo, mockLogger)
	return context.Background(), ctrl, mockRepo, mockLogger, svc
}

func TestCreateMonitor_Success(t *testing.T) {
	tests := []struct {
		name     string
		monitor  models.Monitor
		returnID int64
	}{
		{
			name: "basic monitor",
			monitor: models.Monitor{
				UserID:   1,
				Name:     "Basic Monitor",
				Type:     "http",
				Target:   "https://example.com",
				Timeout:  5,
				Interval: 60,
				IsActive: true,
			},
			returnID: expectedID,
		},
		{
			name: "monitor with JSON specs",
			monitor: models.Monitor{
				UserID:           2,
				Name:             "With JSON",
				Type:             "http",
				Target:           "https://example.org",
				Timeout:          10,
				Interval:         30,
				IsActive:         true,
				RequestSpec:      json.RawMessage(`{"method":"GET"}`),
				ExpectedResponse: json.RawMessage(`{"status":200}`),
			},
			returnID: expectedID,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, ctrl, mockRepo, _, svc := setup(t)
			defer ctrl.Finish()

			mockRepo.EXPECT().
				CreateMonitor(ctx, test.monitor).
				Return(expectedID, nil).
				Times(1)

			id, err := svc.CreateMonitor(ctx, test.monitor)
			assert.NoError(t, err)
			assert.Equal(t, expectedID, id)
		})
	}
}

func TestGetMonitor_Success(t *testing.T) {
	tests := []struct {
		name    string
		id      int64
		monitor *models.Monitor
	}{
		{
			name: "existing monitor",
			id:   expectedID,
			monitor: &models.Monitor{
				ID:     expectedID,
				UserID: 1,
				Name:   "Test Monitor",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, ctrl, mockRepo, _, svc := setup(t)
			defer ctrl.Finish()

			mockRepo.EXPECT().
				GetMonitor(ctx, tt.id).
				Return(tt.monitor, nil).
				Times(1)

			got, err := svc.GetMonitor(ctx, tt.id)
			assert.NoError(t, err)
			assert.Equal(t, tt.monitor, got)
		})
	}
}

func TestGetMonitor_Failure(t *testing.T) {
	tests := []struct {
		name   string
		id     int64
		retErr error
	}{
		{
			name:   "not found",
			id:     expectedID,
			retErr: pgx.ErrNoRows,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, ctrl, mockRepo, _, svc := setup(t)
			defer ctrl.Finish()

			mockRepo.EXPECT().
				GetMonitor(ctx, test.id).
				Return(nil, test.retErr).
				Times(1)

			got, err := svc.GetMonitor(ctx, test.id)
			assert.Error(t, err)
			assert.Nil(t, got)
			assert.Equal(t, test.retErr, err)
		})
	}
}
