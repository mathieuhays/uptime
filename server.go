package uptime

import (
	"github.com/mathieuhays/uptime/internal/handlers"
	"github.com/mathieuhays/uptime/internal/website"
	"html/template"
	"net/http"
)

func NewServer(templ *template.Template, websiteRepository website.Repository) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/{$}", handlers.Home(templ, websiteRepository))
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	return mux
}
