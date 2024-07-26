package uptime

import (
	"database/sql"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mathieuhays/uptime/internal/database"
	"net/http/httptest"
	"strings"
	"testing"
)

const testJwtSecret = "vCAeARsc/ZCXxfVm+2E3Ke0vaNodqVuybZASfxy9Q4IZb+rPWm4ciyIB56uGrdrZMwrxZG7OBWPlyNzFyipGWQ=="

func assertStatus(t testing.TB, response *httptest.ResponseRecorder, want int) {
	t.Helper()
	got := response.Code
	if got != want {
		t.Errorf("did not get expected status. got %v, want %v", got, want)
	}
}

func assertContentType(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("Content-Type") != want {
		t.Errorf("wrong content type. expected %s. got %v", want, response.Result().Header)
	}
}

func assertJSONContentType(t testing.TB, response *httptest.ResponseRecorder) {
	t.Helper()
	assertContentType(t, response, "application/json")
}

func assertBodyContains(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	body := response.Body.String()
	if !strings.Contains(body, want) {
		t.Errorf("expected body to contain %s. got %v", want, body)
	}
}

func assertFieldErrorResponse(t testing.TB, response *httptest.ResponseRecorder, fieldWanted string, messageWanted string) {
	t.Helper()
	var data FieldErrorResponse
	err := json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		t.Fatalf("Couldn't decode Field Error response: %s", err)
	}

	if len(data.Errors) == 0 {
		t.Errorf("Empty field error response. got: %v", response.Body.String())
		return
	}

	foundField := false
	foundMessage := false
	var fieldNames []string

	for _, fieldError := range data.Errors {
		fieldNames = append(fieldNames, fieldError.Field)

		if fieldError.Field == fieldWanted {
			foundField = true

			if fieldError.Message == messageWanted {
				foundMessage = true
			}

			break
		}
	}

	if !foundField {
		t.Errorf("could not find field error. expected: %s. got: %v", fieldWanted, fieldNames)
	}

	if !foundMessage {
		t.Errorf("could not find error messwage. expected: %s. got %v", messageWanted, data)
	}
}

func createTestRouter(t testing.TB) (*Router, *sql.DB, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	config, err := NewApiConfig(database.New(db), testJwtSecret)
	if err != nil {
		t.Fatal(err)
	}

	router, err := NewRouter(config)
	if err != nil {
		t.Fatal(err)
	}

	return router, db, mock
}
