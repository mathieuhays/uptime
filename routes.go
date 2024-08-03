package uptime

import (
	"github.com/mathieuhays/uptime/internal/database"
	"log"
	"net/http"
)

func addRoutes(
	mux *http.ServeMux,
	logger *log.Logger,
	db *database.Queries,
	config *ApiConfig,
) {
	requireAuth := makeRequireAuthMiddleware(db, config.jwtSecret)

	mux.Handle("GET /app/*", http.StripPrefix("/app/", http.FileServer(http.Dir("./static"))))
	mux.Handle("POST /api/v1/users", handleUsersPost(logger, db, config))
	mux.Handle("GET /api/v1/refresh", requireAuth(renderRefresh()))
	mux.Handle("GET /healthz", renderHealth())
	mux.Handle("GET /", http.RedirectHandler("/app/", http.StatusPermanentRedirect))
}
