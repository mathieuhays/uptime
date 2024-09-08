package handlers

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/mathieuhays/uptime/internal/website"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

type websiteRequest struct {
	Name string
	URL  string
}

func (r websiteRequest) Valid(ctx context.Context) (problems map[string]string) {
	problems = map[string]string{}

	if r.Name == "" {
		problems["name"] = "value required"
	}

	if r.URL == "" {
		problems["url"] = "value required"
	} else if !strings.HasPrefix(r.URL, "http") {
		problems["url"] = "invalid URL"
	}

	return problems
}

func Home(templ *template.Template, webRepo website.Repository) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/" {
			writer.WriteHeader(http.StatusNotFound)
			if err := templ.ExecuteTemplate(writer, "404.gohtml", struct{}{}); err != nil {
				log.Printf("error rendering 404: %s\n", err)
			}
			return
		}

		problems := map[string]string{}

		if request.Method == http.MethodPost {
			newWebsiteRequest := websiteRequest{
				Name: request.FormValue("name"),
				URL:  request.FormValue("url"),
			}
			problems = newWebsiteRequest.Valid(request.Context())

			if len(problems) == 0 {
				newWebsite := website.Website{
					ID:        uuid.New(),
					Name:      newWebsiteRequest.Name,
					URL:       newWebsiteRequest.URL,
					CreatedAt: time.Now(),
				}

				if _, err := webRepo.Create(newWebsite); err != nil {
					problems["form"] = fmt.Sprintf("could not create new website. error: %s", err)
				}
			}
		}

		data := struct {
			Websites []website.Website
			Problems map[string]string
		}{
			Websites: []website.Website{},
			Problems: problems,
		}

		webItems, err := webRepo.All()
		if err == nil {
			data.Websites = webItems
		}

		if err := templ.ExecuteTemplate(writer, "index.gohtml", data); err != nil {
			log.Printf("error rendering index: %s\n", err)
		}
	})
}
