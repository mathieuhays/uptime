package website

import (
	"github.com/google/uuid"
	"time"
)

type Website struct {
	ID        uuid.UUID
	Name      string
	URL       string
	CreatedAt time.Time
}

type Repository interface {
	Migrate() error
	Create(website Website) (*Website, error)
	Get(id uuid.UUID) (*Website, error)
	GetByURL(url string) (*Website, error)
	All() ([]Website, error)
}
