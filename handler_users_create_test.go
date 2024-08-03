package uptime

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/mathieuhays/uptime/internal/database"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type StubUserStore struct {
	nextReturn func(id string) (database.User, error)
}

func (s *StubUserStore) Create(ctx context.Context, name, email, encryptedPassword string) (database.User, error) {
	return s.nextReturn("create")
}

func (s *StubUserStore) GetByEmail(ctx context.Context, email string) (database.User, error) {
	return s.nextReturn("email")
}

func (s *StubUserStore) GetByID(ctx context.Context, userID uuid.UUID) (database.User, error) {
	return s.nextReturn("id")
}

type StubSessionStore struct {
	nextReturn func() (database.Session, error)
}

func (s *StubSessionStore) Create(ctx context.Context, userID uuid.UUID) (database.Session, error) {
	return s.nextReturn()
}

func TestHandlerPostUsers(t *testing.T) {
	t.Run("email invalid", func(t *testing.T) {

		server := NewServer(log.Default(), &UserStore{}, &SessionStore{}, &ApiConfig{})
		bodyReader := bytes.NewReader([]byte(`{"name":"John doe","email":"john_doe.com","password":"test123"}`))
		request, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bodyReader)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusBadRequest)
		assertJSONContentType(t, response)
		assertProblemsResponse(t, response, "email", errEmailInvalid.Error())
	})

	t.Run("email already used", func(t *testing.T) {
		userStore := &StubUserStore{nextReturn: func(id string) (database.User, error) {
			if id == "email" {
				return database.User{
					ID:        uuid.New(),
					Name:      "name",
					Email:     "name@email.com",
					Password:  "klsjdflskjdf",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			}

			return database.User{}, errors.New("unexpected state")
		}}

		handler := handleUsersPost(log.Default(), userStore, &StubSessionStore{}, &ApiConfig{})

		bodyReader := bytes.NewReader([]byte(`{"name":"John doe","email":"john@doe.com","password":"test123"}`))
		request, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bodyReader)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusBadRequest)
		assertJSONContentType(t, response)
		assertProblemsResponse(t, response, "email", errEmailAlreadyUsed.Error())
	})

	t.Run("password too short", func(t *testing.T) {
		handler := handleUsersPost(log.Default(), &StubUserStore{}, &StubSessionStore{}, &ApiConfig{})
		bodyReader := bytes.NewReader([]byte(`{"name":"John doe","email":"john@doe.com","password":"test"}`))
		request, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bodyReader)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusBadRequest)
		assertJSONContentType(t, response)
		assertProblemsResponse(t, response, "password", errPasswordTooShort.Error())
	})

	t.Run("password too long", func(t *testing.T) {
		handler := handleUsersPost(log.Default(), &StubUserStore{}, &StubSessionStore{}, &ApiConfig{})
		password := strings.Repeat("é", 72)
		bodyReader := bytes.NewReader([]byte(`{"name":"John doe","email":"john@doe.com","password":"` + password + `"}`))
		request, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bodyReader)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusBadRequest)
		assertJSONContentType(t, response)
		assertProblemsResponse(t, response, "password", errPasswordTooLong.Error())
	})

	t.Run("valid request", func(t *testing.T) {
		userName := "John Doe"
		userEmail := "john@doe.com"
		userPass := "test123"
		refreshToken := "token"

		userStore := &StubUserStore{nextReturn: func(id string) (database.User, error) {
			if id == "create" {
				return database.User{
					ID:        uuid.New(),
					Name:      userName,
					Email:     userEmail,
					Password:  "klsjdflskjdf",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			}

			if id == "email" {
				return database.User{}, sql.ErrNoRows
			}

			return database.User{}, errors.New("unexpected state")
		}}
		sessionStore := &StubSessionStore{nextReturn: func() (database.Session, error) {
			return database.Session{
				ID:           uuid.New(),
				UserID:       uuid.New(),
				RefreshToken: refreshToken,
				ExpireAt:     time.Now(),
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}, nil
		}}

		handler := handleUsersPost(log.Default(), userStore, sessionStore, &ApiConfig{})

		bodyReader := bytes.NewReader([]byte(fmt.
			Sprintf("{\"name\":\"%s\",\"email\":\"%s\",\"password\":\"%s\"}", userName, userEmail, userPass)))

		request, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bodyReader)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusCreated)
		assertJSONContentType(t, response)
		assertBodyContains(t, response, "\"email\":\""+userEmail+"\"")
		assertBodyContains(t, response, "\"refresh_token\":\""+refreshToken+"\"")
	})
}
