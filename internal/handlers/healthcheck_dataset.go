package handlers

import (
	"github.com/google/uuid"
	"github.com/mathieuhays/uptime/internal/healthcheck"
	"net/http"
	"time"
)

type healthcheckDatasetRepo interface {
	GetByWebsiteID(id uuid.UUID, dateRange healthcheck.DateRange) ([]healthcheck.HealthCheck, error)
}

func HealthcheckDataset(repo healthcheckDatasetRepo) http.Handler {
	type point struct {
		X    string `json:"x"`
		Y    int    `json:"y"`
		Code int    `json:"code"`
	}

	type response struct {
		Data []point `json:"data"`
	}

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		id, err := uuid.Parse(request.PathValue("id"))
		if err != nil {
			writer.WriteHeader(http.StatusNotFound)
			return
		}

		var days int
		requestedRange := request.URL.Query().Get("range")
		switch requestedRange {
		case "month":
			days = 30
		case "week":
			days = 7
		default:
			days = 1
		}

		healthChecks, err := repo.GetByWebsiteID(id, healthcheck.DateRange{
			Start: time.Now().Add(time.Hour * 24 * time.Duration(-days)),
			End:   time.Now(),
		})
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		var points []point

		for _, item := range healthChecks {
			points = append(points, point{
				X:    item.CreatedAt.Format(time.DateTime),
				Y:    int(item.ResponseTime.Milliseconds()),
				Code: item.StatusCode,
			})
		}

		if err = encode(writer, request, http.StatusOK, response{Data: points}); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
