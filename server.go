package uptime

import (
	"github.com/mathieuhays/uptime/internal/database"
	"log"
	"net/http"
)

func NewServer(logger *log.Logger, db *database.Queries, config *ApiConfig) http.Handler {
	mux := http.NewServeMux()
	addRoutes(
		mux,
		logger,
		db,
		config,
	)

	// maybe some middleware. logging for example
	// var handler http.Handler = mux
	// handler = someMiddleware(someDependency, handler)
	// return handler

	return mux
}

func renderRefresh() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// todo: implement
		_ = encode(w, r, http.StatusInternalServerError, ErrorResponse{Message: "not implemented"})
	})
}

func renderHealth() http.Handler {
	type response struct {
		Message string `json:"message"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = encode(w, r, http.StatusOK, response{Message: "OK"})
	})
}
