package uptime

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHandlerPostUsers(t *testing.T) {
	router, db, mock := createTestRouter(t)
	defer db.Close()

	t.Run("email invalid", func(t *testing.T) {
		bodyReader := bytes.NewReader([]byte(`{"name":"John doe","email":"john_doe.com","password":"test123"}`))
		request, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bodyReader)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusBadRequest)
		assertJSONContentType(t, response)
		assertProblemsResponse(t, response, "email", errEmailInvalid.Error())

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("expectations not met: %s", err.Error())
		}
	})

	t.Run("email already used", func(t *testing.T) {
		bodyReader := bytes.NewReader([]byte(`{"name":"John doe","email":"john@doe.com","password":"test123"}`))
		request, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bodyReader)
		response := httptest.NewRecorder()

		rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
			AddRow(uuid.New(), "John Doe", "john@doe.com", "some_password_hash", time.Now(), time.Now())
		mock.ExpectQuery("SELECT id, name, email, password, created_at, updated_at FROM users WHERE email").
			WillReturnRows(rows)

		router.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusBadRequest)
		assertJSONContentType(t, response)
		assertProblemsResponse(t, response, "email", errEmailAlreadyUsed.Error())

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("expectations not met: %s", err.Error())
		}
	})

	t.Run("password too short", func(t *testing.T) {
		bodyReader := bytes.NewReader([]byte(`{"name":"John doe","email":"john@doe.com","password":"test"}`))
		request, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bodyReader)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusBadRequest)
		assertJSONContentType(t, response)
		assertProblemsResponse(t, response, "password", errPasswordTooShort.Error())
	})

	t.Run("password too long", func(t *testing.T) {
		password := strings.Repeat("é", 72)
		bodyReader := bytes.NewReader([]byte(`{"name":"John doe","email":"john@doe.com","password":"` + password + `"}`))
		request, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bodyReader)
		response := httptest.NewRecorder()

		mock.ExpectQuery("SELECT id, name, email, password, created_at, updated_at FROM users WHERE email").
			WillReturnError(sql.ErrNoRows)

		router.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusBadRequest)
		assertJSONContentType(t, response)
		assertProblemsResponse(t, response, "password", errPasswordTooLong.Error())

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("expectations not met: %s", err.Error())
		}
	})

	t.Run("valid request", func(t *testing.T) {
		userName := "John Doe"
		userEmail := "john@doe.com"
		userPass := "test123"
		bodyReader := bytes.NewReader([]byte(fmt.
			Sprintf("{\"name\":\"%s\",\"email\":\"%s\",\"password\":\"%s\"}", userName, userEmail, userPass)))

		request, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bodyReader)
		response := httptest.NewRecorder()

		mock.ExpectQuery("SELECT id, name, email, password, created_at, updated_at FROM users WHERE email").
			WillReturnError(sql.ErrNoRows)

		userId := uuid.New()
		userRows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
			AddRow(userId, userName, userEmail, "hash", time.Now(), time.Now())
		mock.ExpectQuery("INSERT INTO users \\(id, name, email, password, created_at, updated_at\\)").
			WillReturnRows(userRows)

		sessionRows := sqlmock.NewRows([]string{"id", "user_id", "refresh_token", "expire_at", "created_at", "updated_at"}).
			AddRow(uuid.New(), userId, "token", time.Now(), time.Now(), time.Now())
		mock.ExpectQuery("INSERT INTO sessions \\(id, user_id, refresh_token, expire_at, created_at, updated_at\\)").
			WillReturnRows(sessionRows)

		router.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusCreated)
		assertJSONContentType(t, response)
		assertBodyContains(t, response, "\"email\":\""+userEmail+"\"")
		assertBodyContains(t, response, "\"refresh_token\":\"token\"")
		assertBodyContains(t, response, "\"access_token\":\"")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("expectations not met: %s", err.Error())
		}
	})
}
