package uptime

import (
	"github.com/mathieuhays/uptime/internal/handlers"
	"github.com/mathieuhays/uptime/internal/website"
	"html/template"
	"net/http"
)

func NewServer(templ *template.Template, websiteRepository website.Repository) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", handlers.Home(templ, websiteRepository))
	mux.Handle("/test", handleTest())

	return mux
}

func handleTest() http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte("uptime!"))
	})
}
