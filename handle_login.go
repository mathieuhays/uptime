package uptime

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/mathieuhays/uptime/internal/database"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
)

func loginCookieRedirect(w http.ResponseWriter, sessionID uuid.UUID) {
	http.SetCookie(w, getCookie(sessionID))
	w.Header().Set("HX-Redirect", "/dashboard")
	w.WriteHeader(http.StatusOK)
}

type loginRequest struct {
	Email    string
	Password string
}

func (l loginRequest) Valid(ctx context.Context) (problems map[string]string) {
	problems = map[string]string{}

	if len(l.Email) == 0 {
		problems["email"] = "Email is required"
	}
	if len(l.Password) == 0 {
		problems["password"] = "Password is required"
	}

	if len(problems) > 0 {
		return problems
	}

	if err := validateEmail(l.Email); err != nil {
		problems["email"] = err.Error()
	}

	return
}

var (
	errInvalidCredentials = errors.New("invalid credentials")
)

func login(loginDetails loginRequest, ctx context.Context, userStore UserStoreInterface, sessionStore SessionStoreInterface) (*database.User, *database.Session, map[string]string, error) {
	problems := loginDetails.Valid(ctx)
	if len(problems) > 0 {
		return nil, nil, problems, errValidation
	}

	user, err := userStore.GetByEmail(ctx, loginDetails.Email)
	if err != nil {
		return nil, nil, nil, errInvalidCredentials
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginDetails.Password)); err != nil {
		return nil, nil, nil, errInvalidCredentials
	}

	session, err := sessionStore.Create(ctx, user.ID)
	if err != nil {
		return nil, nil, nil, errInternalError
	}

	return &user, &session, nil, nil
}

func handleLogin(tmpl *template.Template, userStore UserStoreInterface, sessionStore SessionStoreInterface) http.Handler {
	type formField struct {
		Value string
		Error string
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var fields = map[string]formField{
			"email":    {},
			"password": {},
		}

		var mainError string

		if r.Method == http.MethodPost {
			for fieldName, field := range fields {
				if fieldName != "password" {
					field.Value = r.FormValue(fieldName)
					fields[fieldName] = field
				}
			}

			_, session, problems, err := login(loginRequest{
				Email:    r.FormValue("email"),
				Password: r.FormValue("password"),
			}, r.Context(), userStore, sessionStore)
			if errors.Is(err, errValidation) && len(problems) > 0 {
				for fieldName, problem := range problems {
					if field, ok := fields[fieldName]; ok {
						field.Error = problem
						fields[fieldName] = field
					}
				}
			} else if err != nil {
				mainError = err.Error()
			} else {
				loginCookieRedirect(w, session.ID)
				return
			}
		}

		err := tmpl.ExecuteTemplate(w, "login.gohtml", struct {
			Fields    map[string]formField
			MainError string
			PageTitle string
		}{Fields: fields, MainError: mainError, PageTitle: "Login"})
		if err != nil {
			log.Println(err)
		}
	})
}
