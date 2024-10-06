package uptime

import (
	"github.com/mathieuhays/uptime/internal/handlers"
	"github.com/mathieuhays/uptime/internal/healthcheck"
	"github.com/mathieuhays/uptime/internal/website"
	"html/template"
	"net/http"
)

func NewServer(templ *template.Template, websiteRepository website.Repository, healthCheckRepo healthcheck.Repository) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/{$}", handlers.Home(templ, websiteRepository, healthCheckRepo))
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	mux.Handle("GET /website/{id}", handlers.Website(templ, websiteRepository))
	mux.Handle("/website/{id}/delete", handlers.WebsiteDelete(websiteRepository))
	mux.Handle("GET /website/{id}/healthcheck/dataset", handlers.HealthcheckDataset(healthCheckRepo))

	return mux
}
