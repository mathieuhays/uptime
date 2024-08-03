package uptime

import (
	"log"
	"net/http"
)

func addRoutes(
	mux *http.ServeMux,
	logger *log.Logger,
	userStore UserStoreInterface,
	sessionStore SessionStoreInterface,
	config *ApiConfig,
) {
	requireAuth := makeRequireAuthMiddleware(userStore, config.jwtSecret)

	mux.Handle("/*", handleHome())
	mux.Handle("GET /app/*", http.StripPrefix("/app/", http.FileServer(http.Dir("./static"))))
	mux.Handle("POST /api/v1/users", handleUsersPost(logger, userStore, sessionStore, config))
	mux.Handle("GET /api/v1/refresh", requireAuth(renderRefresh()))
	mux.Handle("GET /healthz", renderHealth())
}

func handleHome() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		http.RedirectHandler("/app/", http.StatusPermanentRedirect).ServeHTTP(w, r)
	})
}
