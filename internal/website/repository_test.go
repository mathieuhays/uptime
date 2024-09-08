package website

import (
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestWebsiteExport(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		id := uuid.New()
		createdAt := time.Now()

		raw := sqliteWebsite{
			uuid:      id.String(),
			name:      "test",
			url:       "https://test.com",
			createdAt: createdAt.Format(time.RFC3339),
		}

		export, err := raw.export()

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if export.ID != id {
			t.Errorf("UUID does not match. expected: %s. got: %s", id, export.ID)
		}

		if export.Name != raw.name {
			t.Errorf("name does not match. expected: %s. got: %s", raw.name, export.Name)
		}

		if export.URL != raw.url {
			t.Errorf("URL does not match. expected: %s. got: %s", raw.url, export.URL)
		}

		if diff := createdAt.Sub(export.CreatedAt).Abs(); diff > time.Second {
			t.Errorf(
				"CreatedAt is above tolerance. diff: %s. expected: %s. got: %s",
				diff,
				createdAt,
				export.CreatedAt)
		}
	})

	t.Run("invalid uuid", func(t *testing.T) {
		raw := sqliteWebsite{
			uuid: "123",
		}

		export, err := raw.export()
		if err == nil {
			t.Errorf("expected a decoding error but got none. id value: %s", export.ID)
		}
	})

	t.Run("invalid createdAt", func(t *testing.T) {
		raw := sqliteWebsite{
			uuid:      uuid.New().String(),
			createdAt: "123456789",
		}

		export, err := raw.export()
		if err == nil {
			t.Errorf("expected a decoding error but got none. date value: %s", export.CreatedAt)
		}
	})
}
