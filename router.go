package uptime

import (
	"github.com/mathieuhays/uptime/internal/database"
	"log"
	"net/http"
)

type Router struct {
	http.Handler
}

func NewRouter(logger *log.Logger, db *database.Queries, config *ApiConfig) (*Router, error) {
	s := new(Router)

	router := http.NewServeMux()
	s.Handler = router

	requireAuth := makeRequireAuthMiddleware(db, config.jwtSecret)

	router.Handle("GET /", http.RedirectHandler("/app/", http.StatusPermanentRedirect))

	router.Handle("GET /app/*", http.StripPrefix("/app/", http.FileServer(http.Dir("./static"))))
	router.Handle("GET /api/v1/health", renderHealth())

	router.Handle("POST /api/v1/users", handleUsersPost(logger, db, config))
	router.Handle("GET /api/v1/refresh", requireAuth(renderRefresh()))

	return s, nil
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
