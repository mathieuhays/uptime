package healthcheck

import (
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestHealthCheckExport(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		id := uuid.New()
		websiteID := uuid.New()
		statusCode := 200
		responseTime := time.Millisecond * 200
		createdAt := time.Now()

		raw := sqliteHealthCheck{
			id:           id.String(),
			websiteID:    websiteID.String(),
			statusCode:   statusCode,
			responseTime: int(responseTime.Milliseconds()),
			createdAt:    int(createdAt.Unix()),
		}

		export, err := raw.export()

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if export.ID != id {
			t.Errorf("UUID does not match for ID. expected: %s. got: %s", id, export.ID)
		}

		if export.WebsiteID != websiteID {
			t.Errorf("UUID does not match for websiteID. expected: %s. got: %s", websiteID, export.WebsiteID)
		}

		if export.StatusCode != statusCode {
			t.Errorf("status code does not match. expected: %d. got: %d", statusCode, export.StatusCode)
		}

		if export.ResponseTime != responseTime {
			t.Errorf("response time does not match. expected: %v. got: %v", responseTime, export.ResponseTime)
		}

		if diff := createdAt.Sub(export.CreatedAt).Abs(); diff > time.Second {
			t.Errorf(
				"CreatedAt is above tolerance. diff: %s. expected: %s. got: %s",
				diff,
				createdAt,
				export.CreatedAt,
			)
		}
	})
}
