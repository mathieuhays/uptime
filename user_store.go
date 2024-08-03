package uptime

import (
	"context"
	"github.com/google/uuid"
	"github.com/mathieuhays/uptime/internal/database"
	"time"
)

type UserStoreInterface interface {
	Create(ctx context.Context, name, email, encryptedPassword string) (database.User, error)
	GetByEmail(ctx context.Context, email string) (database.User, error)
	GetByID(ctx context.Context, userID uuid.UUID) (database.User, error)
}

type UserStore struct {
	db *database.Queries
}

func NewUserStore(db *database.Queries) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) Create(ctx context.Context, name, email, encryptedPassword string) (database.User, error) {
	return s.db.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		Password:  encryptedPassword,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (database.User, error) {
	return s.db.GetUserByEmail(ctx, email)
}

func (s *UserStore) GetByID(ctx context.Context, userID uuid.UUID) (database.User, error) {
	return s.db.GetUserById(ctx, userID)
}
