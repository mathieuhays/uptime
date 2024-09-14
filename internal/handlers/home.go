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

type FormData struct {
	Values websiteRequest
	Errors map[string]string
}

func NewFormData() FormData {
	return FormData{
		Values: websiteRequest{},
		Errors: make(map[string]string),
	}
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
		formData := NewFormData()

		if request.Method == http.MethodPost {
			newWebsiteRequest := websiteRequest{
				Name: request.FormValue("name"),
				URL:  request.FormValue("url"),
			}
			formData.Errors = newWebsiteRequest.Valid(request.Context())
			formData.Values = newWebsiteRequest

			newWebsite := website.Website{
				ID:        uuid.New(),
				Name:      newWebsiteRequest.Name,
				URL:       website.NormalizeURL(newWebsiteRequest.URL),
				CreatedAt: time.Now(),
			}

			if len(formData.Errors) == 0 && urlExists(newWebsite.URL) {
				formData.Errors["url"] = "URL already exists"
			}

			if len(formData.Errors) == 0 {
				if _, err := webRepo.Create(newWebsite); err != nil {
					formData.Errors["form"] = fmt.Sprintln("Something went wrong while creating website. Please try again.")
					log.Printf("website creation error: %s", err)
				} else {
					if request.Header.Get("HX-Request") != "" {
						if err := templ.ExecuteTemplate(writer, "form", NewFormData()); err != nil {
							log.Printf("error rendering index form: %s\n", err)
						}

						if err := templ.ExecuteTemplate(writer, "oob-website-item", newWebsite); err != nil {
							log.Printf("error rendering website-table-item: %s\n", err)
						}
						return
					}

					http.Redirect(writer, request, request.URL.Path, http.StatusFound)
					return
				}
			}

			if len(formData.Errors) > 0 && request.Header.Get("HX-Request") != "" {
				writer.WriteHeader(http.StatusUnprocessableEntity)
				if err := templ.ExecuteTemplate(writer, "form", formData); err != nil {
					log.Printf("error rendering index form: %s\n", err)
				}
				return
			}
		}

		data := struct {
			Websites []website.Website
			FormData FormData
		}{
			Websites: []website.Website{},
			FormData: formData,
		}

		webItems, err := webRepo.All()
		if err == nil {
			data.Websites = webItems
		}

		if err := templ.ExecuteTemplate(writer, "index", data); err != nil {
			log.Printf("error rendering index: %s\n", err)
		}
	})
}
