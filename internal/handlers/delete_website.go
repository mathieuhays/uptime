package handlers

import (
	"github.com/google/uuid"
	"github.com/mathieuhays/uptime/internal/website"
	"net/http"
)

func DeleteWebsite(webRepo website.Repository) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
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

		err = webRepo.Delete(id)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(writer, request, "/", http.StatusFound)
	})
}
