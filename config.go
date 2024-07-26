package uptime

import (
	"github.com/mathieuhays/uptime/internal/database"
	"time"
)

type ApiConfig struct {
	database            *database.Queries
	jwtSecret           string
	sessionDuration     time.Duration
	accessTokenDuration time.Duration
}

func NewApiConfig(db *database.Queries, jwtSecret string) (*ApiConfig, error) {
	return &ApiConfig{
		database:            db,
		jwtSecret:           jwtSecret,
		sessionDuration:     time.Hour * 24 * 2,
		accessTokenDuration: time.Hour,
	}, nil
}
