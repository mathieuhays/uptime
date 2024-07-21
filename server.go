package uptime

import "net/http"

type Server struct {
	http.Handler
}

func NewServer() (*Server, error) {
	s := new(Server)

	router := http.NewServeMux()
	router.Handle("GET /", http.HandlerFunc(renderHomepage))
	router.Handle("GET /api/v1/health", http.HandlerFunc(renderHealth))

	s.Handler = router

	return s, nil
}

func renderHomepage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`<html>
	<head>
		<meta charset="utf-8"/>
		<title>Uptime</title>
	</head>
	<body>
		<h1>Uptime</h1>
	</body>
</html>`))
}

func renderHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
