package website

import (
	"github.com/google/uuid"
	"time"
)

type Website struct {
	ID            uuid.UUID
	Name          string
	URL           string
	LastFetchedAt *time.Time
	CreatedAt     time.Time
}

type Repository interface {
	Create(website Website) (*Website, error)
	Get(id uuid.UUID) (*Website, error)
	GetByURL(url string) (*Website, error)
	Delete(id uuid.UUID) error
	All() ([]Website, error)
	SetAsFetched(id uuid.UUID, date time.Time) error
	GetWebsitesByLastFetched(threshold time.Time, limit int) ([]Website, error)
}
