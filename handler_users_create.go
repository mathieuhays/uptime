package uptime

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/mathieuhays/uptime/internal/database"
	"github.com/mathieuhays/uptime/internal/token"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
	"time"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type FieldErrorResponse struct {
	Errors []FieldError `json:"errors"`
}

var (
	errEmailAlreadyUsed = errors.New("an account is already associated with that email")
)

func (cfg *ApiConfig) handlerPostUsers(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name     string
		Email    string
		Password string
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	email := strings.Trim(payload.Email, " ")
	name := strings.Trim(payload.Name, " ")

	if name == "" || email == "" || payload.Password == "" {
		respondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	var fieldErrors []FieldError

	if err := validateEmail(email); err != nil {
		fieldErrors = append(fieldErrors, FieldError{
			Field:   "email",
			Message: err.Error(),
		})
	} else {
		_, err = cfg.database.GetUserByEmail(r.Context(), email)
		if err == nil {
			fieldErrors = append(fieldErrors, FieldError{
				Field:   "email",
				Message: errEmailAlreadyUsed.Error(),
			})
		} else if !errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusInternalServerError, "something went wrong")
			return
		}
	}

	var encryptedPassword []byte
	var err error

	if err = validatePassword(payload.Password); err != nil {
		fieldErrors = append(fieldErrors, FieldError{
			Field:   "password",
			Message: err.Error(),
		})
	} else {
		encryptedPassword, err = bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			fieldErrors = append(fieldErrors, FieldError{
				Field:   "password",
				Message: errPasswordTooLong.Error(),
			})
		} else if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error encrypting password")
			return
		}
	}

	if len(fieldErrors) > 0 {
		respondWithJSON(w, http.StatusBadRequest, FieldErrorResponse{Errors: fieldErrors})
		return
	}

	encryptedPasswordString := string(encryptedPassword)

	if encryptedPasswordString == "" {
		respondWithError(w, http.StatusInternalServerError, "unexpected state")
		return
	}

	user, err := cfg.database.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		Password:  encryptedPasswordString,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		log.Printf("create_user: user error: %s", err)
		respondWithError(w, http.StatusInternalServerError, "something went wrong while creating user")
		return
	}

	session, err := cfg.database.CreateSession(r.Context(), database.CreateSessionParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		ExpireAt:  time.Now().Add(cfg.sessionDuration).UTC(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		log.Printf("create_user: session error: %s", err)
		respondWithError(w, http.StatusInternalServerError, "something went wrong creating session")
		return
	}

	accessToken, err := token.Generate(user.ID, cfg.jwtSecret, cfg.accessTokenDuration)
	if err != nil {
		log.Printf("creater_user:  access token error: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error authenticating your new user. try to log in.")
		return
	}

	// @TODO add refresh_token and access_token cookie to the responseWriter

	respondWithJSON(w, http.StatusCreated, struct {
		User         User   `json:"user"`
		RefreshToken string `json:"refresh_token"`
		AccessToken  string `json:"access_token"`
	}{
		User:         databaseUserToUser(user),
		RefreshToken: session.RefreshToken,
		AccessToken:  accessToken,
	})
}
