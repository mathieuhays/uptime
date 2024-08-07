package uptime

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/mathieuhays/uptime/internal/auth"
	"github.com/mathieuhays/uptime/internal/token"
	"net/http"
	"time"
)

type ContextKey string

const ContextUserKey ContextKey = "user"

func makeRequireAuthMiddleware(userStore UserStoreInterface, jwtSecret string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessToken, err := auth.GetAuthorization(r.Header)
			if err != nil || accessToken.Name != auth.TypeBearer || accessToken.Value == "" {
				_ = encode(w, r, http.StatusUnauthorized, ErrorResponse{
					Message: "invalid token",
				})
				return
			}

			// @TODO add check for cookies if authorization header is not set

			userID, err := token.Verify(accessToken.Value, jwtSecret)
			if err != nil {
				_ = encode(w, r, http.StatusUnauthorized, ErrorResponse{
					Message: "expired or invalid token",
				})
				return
			}

			user, err := userStore.GetByID(r.Context(), userID)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					_ = encode(w, r, http.StatusUnauthorized, ErrorResponse{Message: "invalid user"})
					return
				}
				_ = encode(w, r, http.StatusInternalServerError, ErrorResponse{
					Message: "error retrieving user",
				})
				return
			}

			ctx := context.WithValue(r.Context(), ContextUserKey, user)

			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func makeRequireLogin(userStore UserStoreInterface, sessionStore SessionStoreInterface) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionCookie, err := r.Cookie("user_session")
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
				return
			}

			sessionID, err := uuid.Parse(sessionCookie.Value)
			if err != nil {
				http.SetCookie(w, cancelCookie())
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
				return
			}

			session, err := sessionStore.Get(r.Context(), sessionID)
			if err != nil {
				http.SetCookie(w, cancelCookie())
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
				return
			}

			user, err := userStore.GetByID(r.Context(), session.UserID)
			if err != nil {
				http.SetCookie(w, cancelCookie())
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
				return
			}

			// refresh cookie
			sessionCookie.Expires = time.Now().Add(time.Hour * 24)
			http.SetCookie(w, sessionCookie)

			c := context.WithValue(r.Context(), ContextUserKey, user)
			h.ServeHTTP(w, r.WithContext(c))
		})
	}
}
