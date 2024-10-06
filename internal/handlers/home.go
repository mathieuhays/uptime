package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/mathieuhays/uptime/internal/healthcheck"
	"github.com/mathieuhays/uptime/internal/website"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

type FormData struct {
	Values struct {
		Name string
		URL  string
	}
	Errors map[string]string
}

func (fd *FormData) Valid(ctx context.Context) (problems map[string]string) {
	problems = map[string]string{}

	if fd.Values.Name == "" {
		problems["name"] = "value required"
	}

	if fd.Values.URL == "" {
		problems["url"] = "value required"
	} else if !strings.HasPrefix(fd.Values.URL, "http") {
		problems["url"] = "invalid URL"
	}

	fd.Errors = problems

	return problems
}

func NewFormData(r *http.Request) FormData {
	var values struct {
		Name string
		URL  string
	}

	if r.Method == http.MethodPost {
		values.Name = r.FormValue("name")
		values.URL = r.FormValue("url")
	}

	return FormData{
		Values: values,
		Errors: make(map[string]string),
	}
}

func makeUrlExists(webRepo website.Repository) func(url string) bool {
	return func(url string) bool {
		_, err := webRepo.GetByURL(url)
		return !errors.Is(err, website.ErrNoRows)
	}
}

func statusCodeToColorCode(statusCode int) string {
	if statusCode < 400 {
		return "success"
	}

	if statusCode < 500 {
		return "warning"
	}

	return "danger"
}

func Home(templ *template.Template, webRepo website.Repository, hcRepo healthcheck.Repository) http.Handler {
	urlExists := makeUrlExists(webRepo)

	type dashboardItem struct {
		Website     website.Website
		HealthCheck *healthcheck.HealthCheck
		ColorCode   string
	}

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		formData := NewFormData(request)

		if request.Method == http.MethodPost {
			formData.Valid(request.Context())

			newWebsite := website.Website{
				ID:        uuid.New(),
				Name:      formData.Values.Name,
				URL:       website.NormalizeURL(formData.Values.URL),
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
						if err := templ.ExecuteTemplate(writer, "form", FormData{}); err != nil {
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
			Items    []dashboardItem
			FormData FormData
		}{
			Items:    []dashboardItem{},
			FormData: formData,
		}

		webItems, err := webRepo.All()
		if err == nil {
			for _, item := range webItems {
				dItem := dashboardItem{
					Website:     item,
					HealthCheck: nil,
					ColorCode:   "",
				}

				hc, err2 := hcRepo.GetLatest(item.ID)
				if err2 == nil {
					dItem.HealthCheck = hc
					dItem.ColorCode = statusCodeToColorCode(hc.StatusCode)
				}

				data.Items = append(data.Items, dItem)
			}
		}

		if err := templ.ExecuteTemplate(writer, "index", data); err != nil {
			log.Printf("error rendering index: %s\n", err)
		}
	})
}
