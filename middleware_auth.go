package uptime

import (
	"database/sql"
	"errors"
	"github.com/mathieuhays/uptime/internal/auth"
	"github.com/mathieuhays/uptime/internal/database"
	"github.com/mathieuhays/uptime/internal/token"
	"net/http"
)

type authenticatedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *ApiConfig) middlewareAuth(handler authenticatedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := auth.GetAuthorization(r.Header)
		if err != nil || accessToken.Name != auth.TypeBearer || accessToken.Value == "" {
			respondWithError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		// @TODO add check for cookies if authorization header is not set

		userId, err := token.Verify(accessToken.Value, cfg.jwtSecret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "expired or invalid token")
			return
		}

		user, err := cfg.database.GetUserById(r.Context(), userId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				respondWithError(w, http.StatusUnauthorized, "invalid user")
				return
			}

			respondWithError(w, http.StatusInternalServerError, "error retrieving user")
			return
		}

		handler(w, r, user)
	}
}
