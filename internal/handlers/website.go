package handlers

import (
	"errors"
	"github.com/google/uuid"
	"github.com/mathieuhays/uptime/internal/healthcheck"
	"github.com/mathieuhays/uptime/internal/website"
	"html/template"
	"log"
	"net/http"
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

		if err = templ.ExecuteTemplate(writer, "website.gohtml", struct {
			Website      website.Website
			HealthChecks []healthcheck.HealthCheck
		}{
			Website:      *w,
			HealthChecks: healthChecks,
		}); err != nil {
			log.Printf("website view. error rendering template: %s", err)
		}
	})
}
