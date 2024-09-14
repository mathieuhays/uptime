package uptime

import (
	"embed"
	"html/template"
	"runtime/debug"
	"slices"
)

//go:embed templates/fragments templates
var templates embed.FS

func TemplateEngine() (*template.Template, error) {
	version := "1"
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		idx := slices.IndexFunc(buildInfo.Settings, func(s debug.BuildSetting) bool {
			return s.Key == "vcs.revision"
		})

		if idx >= 0 {
			version = buildInfo.Settings[idx].Value[:10]
		}
	}

	return template.New("uptime").Funcs(template.FuncMap{
		"assets": func(file string) string {
			return "/static/" + file + "?v=" + version
		},
	}).ParseFS(templates, "templates/*.gohtml", "templates/fragments/*.gohtml")
}
