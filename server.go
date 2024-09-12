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

	mux.Handle("/{$}", handlers.Home(templ, websiteRepository))
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	mux.Handle("GET /website/{id}", handlers.Website(templ, websiteRepository, healthCheckRepo))
	mux.Handle("GET /website/{id}/delete", handlers.DeleteWebsite(websiteRepository))

	return mux
}
