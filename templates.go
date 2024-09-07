package uptime

import (
	"embed"
	"html/template"
)

//go:embed templates/fragments templates
var templates embed.FS

func TemplateEngine() (*template.Template, error) {
	return template.New("uptime").Funcs(template.FuncMap{
		"assets": func(file string) string {
			return "/static/" + file
		},
	}).ParseFS(templates, "templates/*.gohtml", "templates/fragments/*.gohtml")
}
