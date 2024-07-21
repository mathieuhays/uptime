package uptime

import "net/http"

type Router struct {
	http.Handler
}

func NewRouter() (*Router, error) {
	s := new(Router)

	router := http.NewServeMux()

	router.Handle("GET /", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "/app/")
		w.WriteHeader(http.StatusPermanentRedirect)
	}))

	router.Handle("GET /app/*", http.StripPrefix("/app/", http.FileServer(http.Dir("./static"))))
	router.Handle("GET /api/v1/health", http.HandlerFunc(renderHealth))

	s.Handler = router

	return s, nil
}

func renderHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
