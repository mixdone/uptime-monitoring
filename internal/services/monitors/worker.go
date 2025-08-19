package monitors

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/mixdone/uptime-monitoring/internal/models"
	"github.com/mixdone/uptime-monitoring/internal/services/constants"
	"github.com/mixdone/uptime-monitoring/pkg/logger"
	"github.com/mixdone/uptime-monitoring/pkg/message"
)

type worker struct {
	mq     message.MQ
	logger logger.Logger
}

func NewWorker(mq message.MQ, log logger.Logger) Worker {
	return &worker{
		mq:     mq,
		logger: log,
	}
}

func (s *worker) httpHandler(ctx context.Context, data []byte) error {
	var task models.HTTPMonitorTask
	if err := json.Unmarshal(data, &task); err != nil {
		s.logger.WithError(err).Error("failed to unmarshal HTTP monitor task")
		return fmt.Errorf("unmarshal task: %w", err)
	}

	request, err := http.NewRequest(task.Method, task.URL, strings.NewReader(task.RequestBody))
	if err != nil {
		s.logger.WithError(err).Error("failed to create HTTP request")
		return fmt.Errorf("create request error: %w", err)
	}

	for k, v := range task.Headers {
		request.Header.Set(k, v)
	}

	request = request.WithContext(ctx)
	client := &http.Client{Timeout: time.Duration(task.Timeout) * time.Second}

	start := time.Now()
	response, err := client.Do(request)
	duration := time.Since(start).Milliseconds()

	result := models.CheckResult{
		TaskID:    task.ID,
		Duration:  duration,
		CheckedAt: time.Now(),
	}

	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("request failed: %v", err))
		s.logger.WithFields(map[string]any{
			"task_id": task.ID,
			"error":   err,
		}).Error("HTTP request failed")

		return s.mq.Publish("results", toJSON(result))
	}
	defer response.Body.Close()

	if response.StatusCode != task.ExpectStatus {
		result.Errors = append(result.Errors,
			fmt.Sprintf("unexpected status: got %d, want %d",
				response.StatusCode, task.ExpectStatus))
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("failed to read body: %v", err))
	} else {
		body := string(bodyBytes)

		for _, word := range task.ExpectBodyContains {
			if !strings.Contains(body, word) {
				result.Errors = append(result.Errors, fmt.Sprintf("word %q not found in body", word))
			}
		}
	}

	if len(result.Errors) == 0 {
		s.logger.WithFields(map[string]any{
			"task_id":  task.ID,
			"duration": duration,
		}).Info("task succeeded")
	} else {
		s.logger.WithFields(map[string]any{
			"task_id":  task.ID,
			"duration": duration,
			"errors":   result.Errors,
		}).Warn("task failed")
	}

	return s.mq.Publish("results", toJSON(result))
}

func (s *worker) Run(ctx context.Context, monitorType string) {
	handlerFunc := func(data []byte) error {
		switch monitorType {
		case string(MonitorHTTP):
			return s.httpHandler(ctx, data)
		default:
			return errors.New("unknown monitor type")
		}
	}

	backoff := constants.InitialBackoff
	for {
		err := s.mq.Consume(ctx, monitorType, handlerFunc)
		if err == nil {
			s.logger.Infof("%s worker started", monitorType)
			return
		}

		s.logger.WithError(err).Error("failed to start worker, retrying...")
		select {
		case <-time.After(backoff):
			backoff *= 2
			if backoff > constants.MaxBackoff {
				backoff = constants.MaxBackoff
			}
		case <-ctx.Done():
			s.logger.Infof("%s worker stopped", monitorType)
			return
		}
	}
}

func toJSON(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}
