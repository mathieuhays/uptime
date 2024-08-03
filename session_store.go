package uptime

import (
	"context"
	"github.com/google/uuid"
	"github.com/mathieuhays/uptime/internal/database"
	"time"
)

type SessionStoreInterface interface {
	Create(ctx context.Context, userID uuid.UUID) (database.Session, error)
}

type SessionStore struct {
	db     *database.Queries
	config *ApiConfig
}

func NewSessionStore(db *database.Queries, config *ApiConfig) *SessionStore {
	return &SessionStore{db: db, config: config}
}

func (s *SessionStore) Create(ctx context.Context, userID uuid.UUID) (database.Session, error) {
	return s.db.CreateSession(ctx, database.CreateSessionParams{
		ID:        uuid.New(),
		UserID:    userID,
		ExpireAt:  time.Now().Add(s.config.sessionDuration).UTC(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
}
