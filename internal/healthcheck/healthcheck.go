package healthcheck

import (
	"github.com/google/uuid"
	"time"
)

type HealthCheck struct {
	ID           uuid.UUID
	WebsiteID    uuid.UUID
	StatusCode   int
	ResponseTime time.Duration
	CreatedAt    time.Time
}

type Repository interface {
	Create(healthCheck HealthCheck) (*HealthCheck, error)
	Get(id uuid.UUID) (*HealthCheck, error)
	GetByWebsiteID(websiteID uuid.UUID, limit int) ([]HealthCheck, error)
	GetLatest(websiteID uuid.UUID) (*HealthCheck, error)
}
