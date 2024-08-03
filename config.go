package uptime

import (
	"time"
)

type ApiConfig struct {
	jwtSecret           string
	sessionDuration     time.Duration
	accessTokenDuration time.Duration
}

func NewApiConfig(jwtSecret string) (*ApiConfig, error) {
	return &ApiConfig{
		jwtSecret:           jwtSecret,
		sessionDuration:     time.Hour * 24 * 2,
		accessTokenDuration: time.Hour,
	}, nil
}
