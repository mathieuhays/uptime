package handlers

import (
	"errors"
	"github.com/google/uuid"
	"github.com/mathieuhays/uptime/internal/website"
	"net/http"
)

func DeleteWebsite(webRepo website.Repository) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		allowedMethods := map[string]struct{}{
			http.MethodDelete: struct{}{},
			http.MethodGet:    struct{}{},
		}
		if _, ok := allowedMethods[request.Method]; !ok {
			writer.WriteHeader(http.StatusMethodNotAllowed)
		}

		idString := request.PathValue("id")
		if idString == "" {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		id, err := uuid.Parse(idString)
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
