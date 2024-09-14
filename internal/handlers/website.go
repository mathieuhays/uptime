package handlers

import (
	"errors"
	"github.com/google/uuid"
	"github.com/mathieuhays/uptime/internal/healthcheck"
	"github.com/mathieuhays/uptime/internal/website"
	"html/template"
	"log"
	"net/http"
	"time"
)

func Website(templ *template.Template, webRepo website.Repository) http.Handler {
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

		if err = templ.ExecuteTemplate(writer, "website.gohtml", struct {
			Website website.Website
		}{
			Website: *w,
		}); err != nil {
			log.Printf("website view. error rendering template: %s", err)
		}
	})
}

func HealthcheckDataset(hcRepo healthcheck.Repository) http.Handler {
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

		healthChecks, err := hcRepo.GetByWebsiteID(id, healthcheck.DateRange{
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
