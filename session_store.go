package uptime

import (
	"context"
	"github.com/google/uuid"
	"github.com/mathieuhays/uptime/internal/database"
	"net/http"
	"time"
)

const SessionCookie = "user_session"

func getCookie(sessionID uuid.UUID) *http.Cookie {
	return &http.Cookie{
		Name:     SessionCookie,
		Value:    sessionID.String(),
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	}
}

func cancelCookie() *http.Cookie {
	return &http.Cookie{
		Name:     SessionCookie,
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * -1),
		HttpOnly: true,
	}
}

type SessionStoreInterface interface {
	Create(ctx context.Context, userID uuid.UUID) (database.Session, error)
	Get(ctx context.Context, sessionID uuid.UUID) (database.Session, error)
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

func (s *SessionStore) Get(ctx context.Context, sessionID uuid.UUID) (database.Session, error) {
	return s.db.GetSessionByID(ctx, sessionID)
}
