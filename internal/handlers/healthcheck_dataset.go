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
		Data map[string][]point `json:"data"`
	}

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		id, err := uuid.Parse(request.PathValue("id"))
		if err != nil {
			writer.WriteHeader(http.StatusNotFound)
			return
		}

		var days int
		groupInterval := time.Hour * 6
		isSummary := request.URL.Query().Has("summarize")

		requestedRange := request.URL.Query().Get("range")
		switch requestedRange {
		case "month":
			days = 30
		case "week":
			days = 7
			groupInterval = time.Hour * 1
		default:
			days = 1
			groupInterval = time.Minute * 15
		}

		healthChecks, err := repo.GetByWebsiteID(id, healthcheck.DateRange{
			Start: time.Now().Add(time.Hour * 24 * time.Duration(-days)),
			End:   time.Now(),
		})
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !isSummary {
			var points []point

			for _, item := range healthChecks {
				points = append(points, point{
					X:    item.CreatedAt.Format(time.DateTime),
					Y:    int(item.ResponseTime.Milliseconds()),
					Code: item.StatusCode,
				})
			}

			if err = encode(writer, request, http.StatusOK, struct {
				Data []point `json:"data"`
			}{Data: points}); err != nil {
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}

		summarizedItems := healthcheck.SummarizeHealthChecks(healthChecks, groupInterval)

		var averagePoints []point
		var minPoints []point
		var maxPoints []point

		for _, item := range summarizedItems {
			averagePoints = append(averagePoints, point{
				X: item.Date.Format(time.DateTime),
				Y: item.Average,
			})

			minPoints = append(minPoints, point{
				X: item.Date.Format(time.DateTime),
				Y: item.Min,
			})

			maxPoints = append(maxPoints, point{
				X: item.Date.Format(time.DateTime),
				Y: item.Max,
			})
		}

		if err = encode(writer, request, http.StatusOK, struct {
			Data map[string][]point `json:"data"`
		}{Data: map[string][]point{
			"average": averagePoints,
			"min":     minPoints,
			"max":     maxPoints,
		}}); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
