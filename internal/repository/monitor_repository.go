package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mixdone/uptime-monitoring/internal/models"
)

type monitorRepo struct {
	db *pgxpool.Pool
}

func NewMonitorRepo(db *pgxpool.Pool) MonitorsRepository {
	return &monitorRepo{
		db: db,
	}
}

func (r *monitorRepo) CreateMonitor(ctx context.Context, monitor models.Monitor) (int64, error) {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	queryMonitors := `
		INSERT INTO monitors (user_id, name, type, target, timeout, interval, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`

	queryMonitorSpec := `
		INSERT INTO monitor_specs (monitor_id, request, expected_response)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var id int64
	err = tx.QueryRow(ctx, queryMonitors,
		monitor.UserID, monitor.Name, monitor.Type,
		monitor.Target, monitor.Timeout, monitor.Interval,
		monitor.IsActive).Scan(&id)

	if err != nil {
		return 0, err
	}

	var specID int64
	err = tx.QueryRow(ctx, queryMonitorSpec, id, monitor.RequestSpec,
		monitor.ExpectedResponse).Scan(&specID)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *monitorRepo) GetMonitor(ctx context.Context, id int64) (*models.Monitor, error) {

	queryMonitors := ` 
		SELECT m.id, m.user_id, m.name, m.type, m.target, m.timeout, m.interval, 
			m.is_active, m.last_checked_at,
			s.request, s.expected_response
		FROM monitors m
		JOIN monitor_specs s ON m.id = s.monitor_id
		WHERE id = $1
	`

	var monitor models.Monitor
	err := r.db.QueryRow(ctx, queryMonitors, id).Scan(
		&monitor.ID,
		&monitor.UserID,
		&monitor.Name,
		&monitor.Type,
		&monitor.Target,
		&monitor.Timeout,
		&monitor.Interval,
		&monitor.IsActive,
		&monitor.LastCheckedAt,
		&monitor.RequestSpec,
		&monitor.ExpectedResponse)
	if err != nil {
		return nil, err
	}

	return &monitor, nil
}
func (r *monitorRepo) GetAllUserMonitors(ctx context.Context, userID int64) ([]models.Monitor, error) {

	query := ` 
		SELECT 
			m.id, m.user_id, m.name, m.type, m.target, m.timeout, m.interval, 
			m.is_active, m.last_checked_at,
			s.request, s.expected_response
		FROM monitors m
		JOIN monitor_specs s ON s.monitor_id = m.id
		WHERE m.user_id = $1
	`

	var monitors []models.Monitor
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var monitor models.Monitor
		err = rows.Scan(
			&monitor.ID,
			&monitor.UserID,
			&monitor.Name,
			&monitor.Type,
			&monitor.Target,
			&monitor.Timeout,
			&monitor.Interval,
			&monitor.IsActive,
			&monitor.LastCheckedAt,
			&monitor.RequestSpec,
			&monitor.ExpectedResponse)

		if err != nil {
			return nil, err
		}

		monitors = append(monitors, monitor)
	}

	return monitors, nil
}
func (r *monitorRepo) GetAllActiveMonitors(ctx context.Context) ([]models.Monitor, error) {
	query := `
		SELECT m.id, m.user_id, m.name, m.type, m.target, m.timeout, m.interval, 
			m.is_active, m.last_checked_at,
			s.request, s.expected_response
		FROM monitors m
		JOIN monitor_specs s ON m.id = s.monitor_id
		WHERE is_actuve = true 
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var monitors []models.Monitor
	for rows.Next() {
		var monitor models.Monitor
		err = rows.Scan(
			&monitor.ID,
			&monitor.UserID,
			&monitor.Name,
			&monitor.Type,
			&monitor.Target,
			&monitor.Timeout,
			&monitor.Interval,
			&monitor.IsActive,
			&monitor.LastCheckedAt,
			&monitor.RequestSpec,
			&monitor.ExpectedResponse)

		if err != nil {
			return nil, err
		}

		monitors = append(monitors, monitor)
	}

	return monitors, nil
}

func (r *monitorRepo) UpdateMonitor(ctx context.Context, monitor models.Monitor) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	updateMonitorQuery := `
		UPDATE monitors
		SET name = $1, type = $2, target = $3, timeout = $4, interval = $5, is_active = $6
		WHERE id = $7
	`

	_, err = tx.Exec(ctx, updateMonitorQuery,
		monitor.Name,
		monitor.Type,
		monitor.Target,
		monitor.Timeout,
		monitor.Interval,
		monitor.IsActive,
		monitor.ID,
	)
	if err != nil {
		return err
	}

	updateSpecQuery := `
		UPDATE monitor_specs
		SET request = $1, expected_response = $2
		WHERE monitor_id = $3
	`

	_, err = tx.Exec(ctx, updateSpecQuery,
		monitor.RequestSpec,
		monitor.ExpectedResponse,
		monitor.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *monitorRepo) UpdateLastCheckedAt(ctx context.Context, id int64, checkedAt time.Time) error {
	query := `
		UPDATE monitors
		SET last_checked_at = $1
		WHERE id = $2
	`
	_, err := r.db.Exec(ctx, query, checkedAt, id)
	return err
}

func (r *monitorRepo) DeleteMonitor(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM monitors WHERE id = $1`, id)
	return err
}
