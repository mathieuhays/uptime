package uptime

import (
	"github.com/mathieuhays/uptime/internal/database"
	"html/template"
	"log"
	"net/http"
)

func addRoutes(
	mux *http.ServeMux,
	logger *log.Logger,
	templ *template.Template,
	userStore UserStoreInterface,
	sessionStore SessionStoreInterface,
	config *ApiConfig,
) {
	requireAuth := makeRequireAuthMiddleware(userStore, config.jwtSecret)
	requireLogin := makeRequireLogin(userStore, sessionStore)

	mux.Handle("/", requireLogin(handleHome()))
	mux.Handle("GET /dashboard", requireLogin(handleDashboard(logger, templ)))
	mux.Handle("/login", handleLogin(templ, userStore, sessionStore))
	mux.Handle("/logout", handleLogout())
	mux.Handle("/register", handleRegisterHTML(templ, userStore, sessionStore))
	mux.Handle("GET /static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.Handle("POST /api/v1/users", handleUsersPost(userStore, sessionStore, config))
	mux.Handle("GET /api/v1/refresh", requireAuth(renderRefresh()))
	mux.Handle("GET /healthz", renderHealth())
}

func handleHome() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
	})
}

func handleDashboard(logger *log.Logger, tmpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(ContextUserKey).(database.User)
		if !ok {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		err := tmpl.ExecuteTemplate(w, "index.gohtml", struct {
			PageTitle string
			User      database.User
		}{"Homepage", user})
		if err != nil {
			logger.Println(err)
		}
	})
}
