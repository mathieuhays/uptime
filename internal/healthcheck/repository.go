package healthcheck

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type sqliteHealthCheck struct {
	id           string
	websiteID    string
	statusCode   int
	responseTime int
	createdAt    int
}

func (h sqliteHealthCheck) export() (*HealthCheck, error) {
	id, err := uuid.Parse(h.id)
	if err != nil {
		return nil, err
	}

	websiteID, err := uuid.Parse(h.websiteID)
	if err != nil {
		return nil, err
	}

	return &HealthCheck{
		ID:           id,
		WebsiteID:    websiteID,
		StatusCode:   h.statusCode,
		ResponseTime: time.Duration(h.responseTime) * time.Millisecond,
		CreatedAt:    time.Unix(int64(h.createdAt), 0),
	}, nil
}

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) (*SQLiteRepository, error) {
	repo := &SQLiteRepository{db: db}

	_, err := repo.db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, fmt.Errorf("could not enable FOREIGN KEY support: %s", err)
	}

	return repo, nil
}

func (r *SQLiteRepository) Create(healthCheck HealthCheck) (*HealthCheck, error) {
	_, err := r.db.Exec(
		"INSERT INTO healthchecks(id, website_id, status_code, response_time_ms, created_at) VALUES(?, ?, ?, ?, ?);",
		healthCheck.ID.String(),
		healthCheck.WebsiteID.String(),
		healthCheck.StatusCode,
		healthCheck.ResponseTime.Milliseconds(),
		healthCheck.CreatedAt.Unix(),
	)

	if err != nil {
		return nil, err
	}

	return &healthCheck, nil
}

func rowToHealthCheck(row *sql.Row) (*HealthCheck, error) {
	var healthCheck sqliteHealthCheck
	if err := row.Scan(healthCheck.id, healthCheck.websiteID, healthCheck.statusCode, healthCheck.responseTime, healthCheck.createdAt); err != nil {
		return nil, err
	}

	return healthCheck.export()
}

func (r *SQLiteRepository) Get(id uuid.UUID) (*HealthCheck, error) {
	return rowToHealthCheck(r.db.QueryRow(
		"SELECT id, website_id, status_code, response_time_ms, created_at FROM healthchecks WHERE id = ? LIMIT 1;",
		id.String(),
	))
}

func (r *SQLiteRepository) GetByWebsiteID(websiteID uuid.UUID, limit int) ([]HealthCheck, error) {
	rows, err := r.db.Query(
		"SELECT id, website_id, status_code, response_time_ms, created_at FROM healthchecks WHERE website_id = ? ORDER BY created_at DESC LIMIT ?",
		websiteID.String(),
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []HealthCheck
	for rows.Next() {
		var healthCheck sqliteHealthCheck
		if err = rows.Scan(&healthCheck.id, &healthCheck.websiteID, &healthCheck.statusCode, &healthCheck.responseTime, &healthCheck.createdAt); err != nil {
			return nil, err
		}

		h, err := healthCheck.export()
		if err != nil {
			return nil, err
		}

		all = append(all, *h)
	}

	return all, nil
}

func (r *SQLiteRepository) GetLatest(websiteID uuid.UUID) (*HealthCheck, error) {
	return rowToHealthCheck(r.db.QueryRow(
		"SELECT id, website_id, status_code, response_time_ms, created_at FROM healthchecks WHERE website_id = ? ORDER BY created_at DESC LIMIT 1;",
		websiteID.String(),
	))
}
