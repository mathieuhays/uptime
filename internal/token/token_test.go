package token

import (
	"github.com/google/uuid"
	"testing"
	"time"
)

const testingJwtSecret = "iulodofFC4tKO+Gimr8UthVmbEuDb5uW9PYq2c4XGmcU15PpfxEE6625MbNyMcyx7jJ1S71eMpFMQqBBCwsDyQ=="

func TestToken(t *testing.T) {
	t.Run("valid token", func(t *testing.T) {
		id := uuid.New()
		tokenString, err := Generate(id, testingJwtSecret, time.Hour)
		if err != nil {
			t.Fatalf("unexpected generation error: %s", err)
		}

		newId, err := Verify(tokenString, testingJwtSecret)
		if err != nil {
			t.Fatalf("unexpected verification error: %s", err)
		}

		if newId != id {
			t.Errorf("invalid id. expected: %v. got: %v", id, newId)
		}
	})

	t.Run("expired token", func(t *testing.T) {
		tokenString, err := Generate(uuid.New(), testingJwtSecret, time.Millisecond)
		if err != nil {
			t.Fatalf("unexpected generation error: %s", err)
		}

		time.Sleep(time.Millisecond * 5)

		_, err = Verify(tokenString, testingJwtSecret)
		if err == nil {
			t.Error("should have generated an error")
		}
	})
}
