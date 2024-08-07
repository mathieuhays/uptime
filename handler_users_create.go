package uptime

import (
	"context"
	"database/sql"
	"errors"
	"github.com/mathieuhays/uptime/internal/database"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"time"
)

var (
	errEmailAlreadyUsed = errors.New("an account is already associated with that email")
)

type PostUserRequest struct {
	Name     string
	Email    string
	Password string
}

func (r PostUserRequest) Valid(ctx context.Context) (problems map[string]string) {
	problems = map[string]string{}

	if len(r.Name) == 0 {
		problems["name"] = "Name is required"
	}
	if len(r.Email) == 0 {
		problems["email"] = "Email is required"
	}
	if len(r.Password) == 0 {
		problems["password"] = "Password is required"
	}

	if len(problems) > 0 {
		return problems
	}

	if err := validateEmail(r.Email); err != nil {
		problems["email"] = err.Error()
	}

	if err := validatePassword(r.Password); err != nil {
		problems["password"] = err.Error()
	}

	return problems
}

var (
	errValidation            = errors.New("there are validation errors")
	errInternalError         = errors.New("something went wrong")
	errCouldNotCreateUser    = errors.New("could not create user")
	errCouldNotCreateSession = errors.New("could not create session")
)

func signup(userDetails PostUserRequest, ctx context.Context, userStore UserStoreInterface, sessionStore SessionStoreInterface) (*database.User, *database.Session, map[string]string, error) {
	detailsProblems := userDetails.Valid(ctx)
	if len(detailsProblems) > 0 {
		return nil, nil, detailsProblems, errValidation
	}

	problems := map[string]string{}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(userDetails.Password), bcrypt.DefaultCost)
	if errors.Is(err, bcrypt.ErrPasswordTooLong) {
		problems["password"] = errPasswordTooLong.Error()
		return nil, nil, problems, errValidation
	} else if err != nil {
		return nil, nil, problems, errInternalError
	}

	_, err = userStore.GetByEmail(ctx, userDetails.Email)
	if err == nil {
		problems["email"] = errEmailAlreadyUsed.Error()
		return nil, nil, problems, errValidation
	} else if !errors.Is(err, sql.ErrNoRows) {
		return nil, nil, problems, errInternalError
	}

	user, err := userStore.Create(ctx, userDetails.Name, userDetails.Email, string(encryptedPassword))
	if err != nil {
		return nil, nil, problems, errCouldNotCreateUser
	}

	session, err := sessionStore.Create(ctx, user.ID)
	if err != nil {
		return nil, nil, problems, errCouldNotCreateSession
	}

	return &user, &session, problems, nil
}

func handleUsersPost(userStore UserStoreInterface, sessionStore SessionStoreInterface, config *ApiConfig) http.Handler {
	type response struct {
		User         User   `json:"user"`
		RefreshToken string `json:"refresh_token"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userDetails, err := decode[PostUserRequest](r)
		if err != nil {
			_ = encode(w, r, http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		}

		user, session, problems, err := signup(
			userDetails,
			r.Context(),
			userStore,
			sessionStore,
		)
		if errors.Is(err, errValidation) && len(problems) > 0 {
			_ = encode(w, r, http.StatusBadRequest, ErrorResponseWithProblems{Problems: problems})
			return
		} else if err != nil {
			_ = encode(w, r, http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
			return
		}

		// @TODO add refresh_token and access_token cookie to the responseWriter

		_ = encode(w, r, http.StatusCreated, response{
			User:         databaseUserToUser(*user),
			RefreshToken: session.RefreshToken,
		})
	})
}

func handleRegisterHTML(tmpl *template.Template, userStore UserStoreInterface, sessionStore SessionStoreInterface) http.Handler {
	type formField struct {
		Value string
		Error string
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fields := map[string]formField{
			"name":     {},
			"email":    {},
			"password": {},
		}
		var mainError string

		if r.Method == http.MethodPost {
			for fieldName, field := range fields {
				if fieldName == "password" {
					continue
				}

				field.Value = r.FormValue(fieldName)
				fields[fieldName] = field
			}

			_, session, problems, err := signup(
				PostUserRequest{
					Name:     r.FormValue("name"),
					Email:    r.FormValue("email"),
					Password: r.FormValue("password"),
				},
				r.Context(),
				userStore,
				sessionStore,
			)
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
				http.SetCookie(w, &http.Cookie{
					Name:     "user_session",
					Value:    session.ID.String(),
					Path:     "/",
					Expires:  time.Now().Add(24 * time.Hour),
					HttpOnly: true,
					SameSite: 1,
				})
				w.Header().Set("HX-Redirect", "/")
				w.WriteHeader(http.StatusOK)
				return
			}
		} else if r.Method != http.MethodGet {
			http.NotFound(w, r)
			return
		}

		err := tmpl.ExecuteTemplate(w, "register.gohtml", struct {
			Fields    map[string]formField
			MainError string
			PageTitle string
		}{Fields: fields, MainError: mainError, PageTitle: "Register"})
		if err != nil {
			log.Println(err)
		}
	})
}
