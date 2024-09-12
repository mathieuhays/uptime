package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func ErrorPageFromStatus(templ *template.Template, writer http.ResponseWriter, statusCode int) {
	writer.WriteHeader(statusCode)

	if err := templ.ExecuteTemplate(writer, "error.gohtml", struct {
		Code  int
		Title string
	}{
		Code:  statusCode,
		Title: http.StatusText(statusCode),
	}); err != nil {
		log.Printf("error rendering error view: %s\n", err)
	}
}
