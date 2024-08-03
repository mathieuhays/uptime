package uptime

import (
	"context"
	"database/sql"
	"errors"
	"github.com/mathieuhays/uptime/internal/token"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
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

func handleUsersPost(logger *log.Logger, userStore *UserStore, sessionStore *SessionStore, config *ApiConfig) http.Handler {
	type response struct {
		User         User   `json:"user"`
		RefreshToken string `json:"refresh_token"`
		AccessToken  string `json:"access_token"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userDetails, problems, err := decodeValid[PostUserRequest](r)
		if len(problems) > 0 {
			_ = encode(w, r, http.StatusBadRequest, ErrorResponseWithProblems{Problems: problems})
			return
		} else if err != nil {
			_ = encode(w, r, http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		}

		_, err = userStore.GetByEmail(r.Context(), userDetails.Email)
		if err == nil {
			_ = encode(w, r, http.StatusBadRequest, ErrorResponseWithProblems{
				Problems: map[string]string{"email": errEmailAlreadyUsed.Error()},
			})
			return
		} else if !errors.Is(err, sql.ErrNoRows) {
			logger.Printf("handle users post, check existing user err: %s", err)
			_ = encode(w, r, http.StatusInternalServerError, ErrorResponse{Message: "Something went wrong"})
			return
		}

		if err = validatePassword(userDetails.Password); err != nil {
			_ = encode(w, r, http.StatusBadRequest, ErrorResponseWithProblems{
				Problems: map[string]string{"password": err.Error()},
			})
			return
		}

		encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(userDetails.Password), bcrypt.DefaultCost)
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			_ = encode(w, r, http.StatusBadRequest, ErrorResponseWithProblems{
				Problems: map[string]string{"password": errPasswordTooLong.Error()},
			})
			return
		} else if err != nil {
			logger.Printf("post users: password encoding err: %s", err)
			_ = encode(w, r, http.StatusInternalServerError, ErrorResponse{Message: "Something went wrong"})
			return
		}

		user, err := userStore.Create(r.Context(), userDetails.Name, userDetails.Email, string(encryptedPassword))
		if err != nil {
			logger.Printf("post users: CreateUser err: %s", err)
			_ = encode(w, r, http.StatusInternalServerError, ErrorResponse{
				Message: "Your user could not be created, please try again later",
			})
			return
		}

		session, err := sessionStore.Create(r.Context(), user.ID)
		if err != nil {
			logger.Printf("create_user: session error: %s", err)
			_ = encode(w, r, http.StatusInternalServerError, ErrorResponse{
				Message: "something went wrong creating session",
			})
			return
		}

		accessToken, err := token.Generate(user.ID, config.jwtSecret, config.accessTokenDuration)
		if err != nil {
			logger.Printf("creater_user: access token error: %s", err)
			_ = encode(w, r, http.StatusInternalServerError, ErrorResponse{
				Message: "error authenticating your new user. try to log in",
			})
			return
		}

		// @TODO add refresh_token and access_token cookie to the responseWriter

		_ = encode(w, r, http.StatusCreated, response{
			User:         databaseUserToUser(user),
			RefreshToken: session.RefreshToken,
			AccessToken:  accessToken,
		})
	})
}
