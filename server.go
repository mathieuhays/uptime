package uptime

import (
	"html/template"
	"log"
	"net/http"
)

func NewServer(templ *template.Template) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", handleHome(templ))
	mux.Handle("/test", handleTest())

	return mux
}

func handleHome(templ *template.Template) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("rendering home\n")

		if err := templ.ExecuteTemplate(writer, "index.gohtml", struct{}{}); err != nil {
			log.Printf("error rendering index: %s\n", err)
		}
	})
}

func handleTest() http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte("uptime!"))
	})
}
