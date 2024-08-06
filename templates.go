package uptime

import (
	"embed"
	"html/template"
)

//go:embed templates/fragments templates
var pageTemplates embed.FS

func GetTemplateEngine() (*template.Template, error) {
	return template.New("uptime").Funcs(template.FuncMap{
		"assets": func(file string) string {
			return "/static/" + file
		},
	}).ParseFS(pageTemplates, "templates/*.gohtml", "templates/fragments/*.gohtml")
}
