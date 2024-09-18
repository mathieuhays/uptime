package handlers

import (
	"errors"
	"github.com/google/uuid"
	"github.com/mathieuhays/uptime/internal/website"
	"net/http"
)

type websiteDeleteRepo interface {
	Get(id uuid.UUID) (*website.Website, error)
	Delete(id uuid.UUID) error
}

func WebsiteDelete(webRepo websiteDeleteRepo) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		allowedMethods := map[string]struct{}{
			http.MethodDelete: {},
			http.MethodGet:    {},
		}
		if _, ok := allowedMethods[request.Method]; !ok {
			writer.WriteHeader(http.StatusMethodNotAllowed)
		}

		id, err := uuid.Parse(request.PathValue("id"))
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err = webRepo.Get(id)
		if errors.Is(err, website.ErrNoRows) {
			writer.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = webRepo.Delete(id)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		if request.Header.Get("HX-Request") != "" {
			writer.WriteHeader(http.StatusOK)
			return
		}

		http.Redirect(writer, request, "/", http.StatusFound)
	})
}
