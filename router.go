package uptime

import "net/http"

type Router struct {
	http.Handler
	config *ApiConfig
}

func NewRouter(config *ApiConfig) (*Router, error) {
	s := new(Router)

	router := http.NewServeMux()
	s.Handler = router
	s.config = config

	router.Handle("GET /", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "/app/")
		w.WriteHeader(http.StatusPermanentRedirect)
	}))

	router.Handle("GET /app/*", http.StripPrefix("/app/", http.FileServer(http.Dir("./static"))))
	router.Handle("GET /api/v1/health", http.HandlerFunc(renderHealth))

	router.Handle("POST /api/v1/users", http.HandlerFunc(config.handlerPostUsers))

	return s, nil
}

func renderHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
