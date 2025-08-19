package monitors

import (
	"container/heap"
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/mixdone/uptime-monitoring/internal/models"
	"github.com/mixdone/uptime-monitoring/internal/models/errs"
	"github.com/mixdone/uptime-monitoring/internal/repository"
	"github.com/mixdone/uptime-monitoring/internal/services/constants"
	"github.com/mixdone/uptime-monitoring/pkg/logger"
	"github.com/mixdone/uptime-monitoring/pkg/message"
)

type MonitorType string

const (
	MonitorHTTP MonitorType = "http"
	ResultChan  string      = "result"
)

type ScheduledTask struct {
	MonitorID   int64
	Kind        MonitorType
	NextCheckAt time.Time
	Interval    time.Duration
	index       int
}

type manager struct {
	mq     message.MQ
	logger logger.Logger
	mutex  sync.Mutex
	h      taskHeap
	inHeap map[int64]bool
	repo   repository.MonitorsRepository
}

func NewManager(mq message.MQ, log logger.Logger, repo repository.MonitorsRepository) Manager {
	return &manager{
		mq:     mq,
		logger: log,
		h:      taskHeap{},
		inHeap: make(map[int64]bool),
		repo:   repo,
	}
}

func (m *manager) AddTask(t *ScheduledTask) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.inHeap[t.MonitorID] {
		return
	}

	heap.Push(&m.h, t)
	m.inHeap[t.MonitorID] = true
}

func (m *manager) DequeueTask(id int64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.inHeap[id] = false
}

func (m *manager) Run(ctx context.Context) {

	monitors, err := m.repo.GetAllActiveMonitors(ctx)

	if err == nil {
		for _, monitor := range monitors {
			m.AddTask(&ScheduledTask{
				MonitorID:   monitor.ID,
				Kind:        MonitorType(monitor.Type),
				NextCheckAt: time.Now(),
				Interval:    time.Duration(monitor.Interval),
			})
		}
	} else {
		m.logger.WithError(err).Warn("failed to fetch active monitors from db")
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			m.logger.Info("manager Stoped")
			return
		case <-ticker.C:
			m.dispatchTasks()
		}
	}
}

func (m *manager) dispatchTasks() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for m.h.Len() > 0 {
		if time.Now().Before(m.h[0].NextCheckAt) {
			return
		}

		task := heap.Pop(&m.h).(*ScheduledTask)

		if !m.inHeap[task.MonitorID] {
			continue
		}

		go func() {

		}()

		if err := m.publishTask(task); err != nil {
			m.logger.WithError(err).Error("failed to publish task")
		}

		task.NextCheckAt = time.Now().Add(time.Duration(task.Interval) * time.Second)
		heap.Push(&m.h, task)
	}
}

func (m *manager) publishTask(task *ScheduledTask) error {
	monitor, err := m.repo.GetMonitor(context.Background(), task.MonitorID)
	if err != nil {
		return err
	}

	var publishTask interface{}
	switch task.Kind {
	case MonitorHTTP:
		publishTask, err = monitor.ToHTTPMonitorTask()
	default:
		return errs.ErrWrongType
	}

	if err != nil {
		m.logger.WithError(err).Error("failed to convert monitor")
		return err
	}

	data, err := json.Marshal(publishTask)
	if err != nil {
		return err
	}

	return m.mq.Publish(string(task.Kind), data)
}

func (m *manager) StartResultHandler(ctx context.Context) {

	backoff := constants.InitialBackoff
	for {
		err := m.mq.Consume(ctx, ResultChan, m.HandleResult)

		if err == nil {
			m.logger.Info("result handler started")
			return
		}

		m.logger.WithError(err).Error("failed to start result handler")

		select {
		case <-time.After(backoff):
			backoff *= 2
			if backoff > constants.MaxBackoff {
				backoff = constants.MaxBackoff
			}
		case <-ctx.Done():
			m.logger.Infof("stopping result handler")
			return
		}
	}
}

func (m *manager) HandleResult(data []byte) error {
	var result models.CheckResult
	if err := json.Unmarshal(data, &result); err != nil {
		m.logger.WithError(err).Error("failed to unmarshal result")
		return err
	}

	if len(result.Errors) > 0 {
		m.logger.WithFields(map[string]any{
			"task_id":  result.TaskID,
			"duration": result.Duration,
			"errors":   result.Errors,
		}).Warn("task failed")
	}

	m.logger.WithFields(map[string]any{
		"task_id":  result.TaskID,
		"duration": result.Duration,
	}).Info("task succeeded")

	return nil
}
