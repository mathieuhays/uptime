package handlers

import (
	"github.com/google/uuid"
	"github.com/mathieuhays/uptime/internal/asserts"
	"github.com/mathieuhays/uptime/internal/website"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type stubWebsiteDeleteRepo struct {
	get    func(id uuid.UUID) (*website.Website, error)
	delete func(id uuid.UUID) error
}

func (s stubWebsiteDeleteRepo) Get(id uuid.UUID) (*website.Website, error) {
	if s.get != nil {
		return s.get(id)
	}

	return &website.Website{
		ID:            id,
		Name:          "test",
		URL:           "https://example.com",
		LastFetchedAt: nil,
		CreatedAt:     time.Now(),
	}, nil
}
func (s stubWebsiteDeleteRepo) Delete(id uuid.UUID) error {
	if s.delete != nil {
		return s.delete(id)
	}

	return nil
}

func TestDeleteWebsite(t *testing.T) {
	t.Run("block unauthorized methods", func(t *testing.T) {
		handler := WebsiteDelete(stubWebsiteDeleteRepo{})
		methods := []string{
			http.MethodPost,
			http.MethodPut,
		}

		for _, method := range methods {
			request, err := http.NewRequest(method, "/", nil)
			request.SetPathValue("id", uuid.New().String())

			if err != nil {
				t.Errorf("unexpected error when generating request for %s method: %s", method, err)
				continue
			}

			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)

			asserts.StatusCode(t, response, http.StatusMethodNotAllowed)
		}
	})

	t.Run("handle invalid ids", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodDelete, "/", nil)
		if err != nil {
			t.Fatalf("unexpected error when building request: %s", err)
		}

		request.SetPathValue("id", "lskjdf")
		response := httptest.NewRecorder()
		WebsiteDelete(stubWebsiteDeleteRepo{}).ServeHTTP(response, request)

		asserts.StatusCode(t, response, http.StatusBadRequest)
	})

	t.Run("redirected when accessed directly", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Fatalf("unexpected error when building request: %s", err)
		}

		id := uuid.New()
		request.SetPathValue("id", id.String())
		response := httptest.NewRecorder()
		WebsiteDelete(stubWebsiteDeleteRepo{}).ServeHTTP(response, request)

		asserts.StatusCode(t, response, http.StatusFound)
	})

	t.Run("simple OK when using HTMX", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodDelete, "/", nil)
		if err != nil {
			t.Fatalf("unexpected error when building request: %s", err)
		}

		id := uuid.New()
		request.SetPathValue("id", id.String())
		request.Header.Set("HX-Request", "true")
		response := httptest.NewRecorder()
		WebsiteDelete(stubWebsiteDeleteRepo{}).ServeHTTP(response, request)

		asserts.StatusCode(t, response, http.StatusOK)
	})

	t.Run("404 when resource is already deleted", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Fatalf("unexpected error when building request: %s", err)
		}

		id := uuid.New()
		request.SetPathValue("id", id.String())
		response := httptest.NewRecorder()

		WebsiteDelete(stubWebsiteDeleteRepo{
			get: func(id uuid.UUID) (*website.Website, error) {
				return nil, website.ErrNoRows
			},
		}).ServeHTTP(response, request)

		asserts.StatusCode(t, response, http.StatusNotFound)
	})
}
