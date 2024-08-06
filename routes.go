package uptime

import (
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

	mux.Handle("/*", handleHome(logger, templ))
	mux.Handle("/login", handleLogin(logger, templ))
	mux.Handle("GET /static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.Handle("POST /api/v1/users", handleUsersPost(logger, userStore, sessionStore, config))
	mux.Handle("GET /api/v1/refresh", requireAuth(renderRefresh()))
	mux.Handle("GET /healthz", renderHealth())
}

func handleHome(logger *log.Logger, tmpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		err := tmpl.ExecuteTemplate(w, "index.gohtml", struct {
			PageTitle string
		}{"Homepage"})
		if err != nil {
			logger.Println(err)
		}
	})
}

func handleLogin(logger *log.Logger, tmpl *template.Template) http.Handler {
	type formField struct {
		Value string
		Error string
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var fields = map[string]formField{}

		if r.Method == http.MethodPost {
			// handle login action
			email := r.FormValue("email")
			emailField := formField{Value: email}

			if err := validateEmail(email); err != nil {
				emailField.Error = err.Error()
			}

			fields["email"] = emailField

			// on success redirect to dashboard or something
		}

		err := tmpl.ExecuteTemplate(w, "login.gohtml", struct {
			Fields    map[string]formField
			PageTitle string
		}{Fields: fields, PageTitle: "Login"})
		if err != nil {
			log.Println(err)
		}
	})
}
