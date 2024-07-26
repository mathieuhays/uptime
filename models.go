package uptime

import (
	"github.com/google/uuid"
	"github.com/mathieuhays/uptime/internal/database"
)

type User struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

func databaseUserToUser(user database.User) User {
	return User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}

type Session struct {
	ID           uuid.UUID `json:"id"`
	RefreshToken string    `json:"refresh_token"`
}

func databaseSessionToSession(session database.Session) Session {
	return Session{
		ID:           session.ID,
		RefreshToken: session.RefreshToken,
	}
}
