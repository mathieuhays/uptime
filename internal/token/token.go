package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type customClaims struct {
	jwt.RegisteredClaims
}

func Generate(userId uuid.UUID, secret string, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "uptime",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration).UTC()),
		Subject:   userId.String(),
	})
	return token.SignedString([]byte(secret))
}

func Verify(accessToken string, secret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(accessToken, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	claims, ok := token.Claims.(*customClaims)
	if !ok {
		return uuid.UUID{}, errors.New("could not type cast claims")
	}

	return uuid.Parse(claims.Subject)
}
