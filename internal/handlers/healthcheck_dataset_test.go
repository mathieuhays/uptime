package handlers

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mathieuhays/uptime/internal/asserts"
	"github.com/mathieuhays/uptime/internal/healthcheck"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type healthCheckDatasetResponse struct {
	Data []struct {
		X    string
		Y    int
		Code int
	}
}

type stubHealthCheckDatasetRepo struct {
	get func(id uuid.UUID, dateRange healthcheck.DateRange) ([]healthcheck.HealthCheck, error)
}

func (s stubHealthCheckDatasetRepo) GetByWebsiteID(id uuid.UUID, dateRange healthcheck.DateRange) ([]healthcheck.HealthCheck, error) {
	if s.get != nil {
		return s.get(id, dateRange)
	}

	return []healthcheck.HealthCheck{}, nil
}

func TestHealthcheckDataset(t *testing.T) {
	testCases := []struct {
		rangeName string
		duration  time.Duration
	}{
		{
			rangeName: "month",
			duration:  time.Hour * 24 * 30,
		},
		{
			rangeName: "week",
			duration:  time.Hour * 24 * 7,
		},
		{
			rangeName: "day",
			duration:  time.Hour * 24,
		},
		{
			rangeName: "invalid",
			duration:  time.Hour * 24,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s range", tc.rangeName), func(t *testing.T) {
			request, err := http.NewRequest(http.MethodGet, "/?range="+tc.rangeName, nil)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			request.SetPathValue("id", uuid.New().String())
			threshold := time.Now().Add(-tc.duration)
			var requestedDateRange *healthcheck.DateRange
			response := httptest.NewRecorder()

			HealthcheckDataset(stubHealthCheckDatasetRepo{get: func(id uuid.UUID, dateRange healthcheck.DateRange) ([]healthcheck.HealthCheck, error) {
				requestedDateRange = &dateRange
				return []healthcheck.HealthCheck{}, nil
			}}).ServeHTTP(response, request)

			asserts.StatusCode(t, response, http.StatusOK)

			if requestedDateRange == nil {
				t.Fatalf("no date range requested. one is expected even in invalid cases")
			}

			asserts.SimilarTime(t, requestedDateRange.Start, threshold)
		})
	}

	t.Run("invalid id", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		request.SetPathValue("id", "lkjsdlf")

		response := httptest.NewRecorder()
		HealthcheckDataset(stubHealthCheckDatasetRepo{}).ServeHTTP(response, request)

		asserts.StatusCode(t, response, http.StatusNotFound)
	})

	t.Run("empty results", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		request.SetPathValue("id", uuid.New().String())
		response := httptest.NewRecorder()

		HealthcheckDataset(stubHealthCheckDatasetRepo{}).ServeHTTP(response, request)

		asserts.StatusCode(t, response, http.StatusOK)

		res, err := decode[healthCheckDatasetResponse](response.Body)
		if err != nil {
			t.Fatalf("unexpected error decoding JSON response: %s", err)
		}

		if len(res.Data) != 0 {
			t.Errorf("unexpected data. should be empty. got: %v", res.Data)
		}
	})

	t.Run("results", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		request.SetPathValue("id", uuid.New().String())
		response := httptest.NewRecorder()

		expectedCode := 200
		expectedY := 300
		expectedTime := time.Now()
		expectedTimeString := expectedTime.Format(time.DateTime)

		HealthcheckDataset(stubHealthCheckDatasetRepo{get: func(id uuid.UUID, dateRange healthcheck.DateRange) ([]healthcheck.HealthCheck, error) {
			return []healthcheck.HealthCheck{
				{
					ID:           uuid.New(),
					WebsiteID:    id,
					StatusCode:   expectedCode,
					ResponseTime: time.Duration(expectedY) * time.Millisecond,
					CreatedAt:    expectedTime,
				},
			}, nil
		}}).ServeHTTP(response, request)

		asserts.StatusCode(t, response, http.StatusOK)

		res, err := decode[healthCheckDatasetResponse](response.Body)
		if err != nil {
			t.Fatalf("unexpected error decoding JSON response: %s", err)
		}

		if len(res.Data) != 1 {
			t.Fatalf("unexpected data. should have 1 item. got: %v", res.Data)
		}

		if res.Data[0].Code != expectedCode {
			t.Errorf("wrong code. expected: %d. got: %d", expectedCode, res.Data[0].Code)
		}

		if res.Data[0].Y != expectedY {
			t.Errorf("wrong response time (Y). expected: %d. got: %d", expectedY, res.Data[0].Y)
		}

		if res.Data[0].X != expectedTimeString {
			t.Errorf("wrong time (X). expected: %s. got: %s", expectedTimeString, res.Data[0].X)
		}
	})
}
