package uptime

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServerHomepage(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("unexpected status code. expected: 200. got: %d", response.Code)
	}

	if response.Header().Get("Content-Type") != "text/html" {
		t.Errorf("wrong content type. expected: text/html. got: %s", response.Header().Get("Content-Type"))
	}
}

func TestServerApiHealth(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	request, _ := http.NewRequest(http.MethodGet, "/api/v1/health", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("unexpected status code. expected: 200. got: %d", response.Code)
	}

	if response.Header().Get("Content-Type") != "application/json" {
		t.Errorf("wrong content type. expected: application/json. got: %s", response.Header().Get("Content-Type"))
	}
}
