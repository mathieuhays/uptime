package handlers

import (
	"context"
	"errors"
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

func makeUrlExists(webRepo website.Repository) func(url string) bool {
	return func(url string) bool {
		_, err := webRepo.GetByURL(url)
		return !errors.Is(err, website.ErrNoRows)
	}
}

func Home(templ *template.Template, webRepo website.Repository) http.Handler {
	urlExists := makeUrlExists(webRepo)

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		problems := map[string]string{}
		var newWebsiteRequest websiteRequest

		if request.Method == http.MethodPost {
			newWebsiteRequest = websiteRequest{
				Name: request.FormValue("name"),
				URL:  request.FormValue("url"),
			}
			problems = newWebsiteRequest.Valid(request.Context())

			newWebsite := website.Website{
				ID:        uuid.New(),
				Name:      newWebsiteRequest.Name,
				URL:       website.NormalizeURL(newWebsiteRequest.URL),
				CreatedAt: time.Now(),
			}

			if len(problems) == 0 && urlExists(newWebsite.URL) {
				problems["url"] = "URL already exists"
			}

			if len(problems) == 0 {
				if _, err := webRepo.Create(newWebsite); err != nil {
					problems["form"] = fmt.Sprintln("Something went wrong while creating website. Please try again.")
					log.Printf("website creation error: %s", err)
				} else {
					http.Redirect(writer, request, request.URL.Path, http.StatusFound)
					return
				}
			}
		}

		data := struct {
			Websites []website.Website
			Problems map[string]string
			Values   websiteRequest
		}{
			Websites: []website.Website{},
			Problems: problems,
			Values:   newWebsiteRequest,
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
