package uptime

import (
	"html/template"
	"log"
	"net/http"
)

func NewServer(templ *template.Template) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", handleHome(templ))

	return mux
}

func handleHome(templ *template.Template) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if err := templ.ExecuteTemplate(writer, "index.gohtml", struct{}{}); err != nil {
			log.Printf("error rendering index: %s\n", err)
		}
	})
}
