package uptime

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterRoot(t *testing.T) {
	router := NewServer(log.Default(), nil, &ApiConfig{})
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assertStatus(t, response, http.StatusPermanentRedirect)

	location := response.Header().Get("Location")
	expectedLocation := "/app/"
	if location != expectedLocation {
		t.Errorf("Invalid redirect. expected: %q. got: %q", expectedLocation, location)
	}
}

func TestRouterAppHomepage(t *testing.T) {
	router := NewServer(log.Default(), nil, &ApiConfig{})
	request, _ := http.NewRequest(http.MethodGet, "/app/", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assertStatus(t, response, http.StatusOK)
	assertContentType(t, response, "text/html; charset=utf-8")
}

func TestRouterHealth(t *testing.T) {
	router := NewServer(log.Default(), nil, &ApiConfig{})
	request, _ := http.NewRequest(http.MethodGet, "/healthz", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("unexpected status code. expected: 200. got: %d", response.Code)
	}

	if response.Header().Get("Content-Type") != "application/json" {
		t.Errorf("wrong content type. expected: application/json. got: %s", response.Header().Get("Content-Type"))
	}
}
