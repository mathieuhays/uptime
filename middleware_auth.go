package uptime

import (
	"context"
	"database/sql"
	"errors"
	"github.com/mathieuhays/uptime/internal/auth"
	"github.com/mathieuhays/uptime/internal/database"
	"github.com/mathieuhays/uptime/internal/token"
	"net/http"
)

type ContextKey string

const ContextUserKey ContextKey = "user"

func makeRequireAuthMiddleware(db *database.Queries, jwtSecret string) func(h http.Handler) http.Handler {
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

			userId, err := token.Verify(accessToken.Value, jwtSecret)
			if err != nil {
				_ = encode(w, r, http.StatusUnauthorized, ErrorResponse{
					Message: "expired or invalid token",
				})
				return
			}

			user, err := db.GetUserById(r.Context(), userId)
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
