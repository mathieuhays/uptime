package handlers

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/mathieuhays/uptime/internal/healthcheck"
	"github.com/mathieuhays/uptime/internal/website"
	"html/template"
	"log"
	"net/http"
	"time"
)

func Website(templ *template.Template, webRepo website.Repository, hcRepo healthcheck.Repository) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		id, err := uuid.Parse(request.PathValue("id"))
		if err != nil {
			ErrorPageFromStatus(templ, writer, http.StatusBadRequest)
			return
		}

		w, err := webRepo.Get(id)
		if errors.Is(err, website.ErrNoRows) {
			ErrorPageFromStatus(templ, writer, http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("website single. query error: %s\n", err)
			ErrorPageFromStatus(templ, writer, http.StatusInternalServerError)
			return
		}

		healthChecks, err := hcRepo.GetByWebsiteID(w.ID, 50)
		if err != nil {
			ErrorPageFromStatus(templ, writer, http.StatusInternalServerError)
			return
		}

		var healthCheckChartData string
		healthCheckChartDataBytes, err := json.Marshal(healthChecksToChartData(healthChecks))
		if err == nil {
			healthCheckChartData = string(healthCheckChartDataBytes)
		}

		if err = templ.ExecuteTemplate(writer, "website.gohtml", struct {
			Website              website.Website
			HealthChecks         []healthcheck.HealthCheck
			HealthCheckChartData string
		}{
			Website:              *w,
			HealthChecks:         healthChecks,
			HealthCheckChartData: healthCheckChartData,
		}); err != nil {
			log.Printf("website view. error rendering template: %s", err)
		}
	})
}

type point struct {
	X string `json:"x"`
	Y int    `json:"y"`
}

func healthChecksToChartData(items []healthcheck.HealthCheck) []point {
	points := []point{}

	for _, item := range items {
		points = append(points, point{
			X: item.CreatedAt.Format(time.DateTime),
			Y: int(item.ResponseTime.Milliseconds()),
		})
	}

	return points
}
