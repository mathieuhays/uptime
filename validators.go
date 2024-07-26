package uptime

import (
	"errors"
	"net/mail"
	"strings"
)

const passwordMinLength = 6

var (
	errEmailInvalid     = errors.New("invalid email")
	errPasswordTooShort = errors.New("password is too short")
	errPasswordTooLong  = errors.New("password is too long")
)

func validateEmail(email string) error {
	if email == "" {
		return errEmailInvalid
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return errEmailInvalid
	}

	// missing tld
	if !strings.Contains(parts[1], ".") {
		return errEmailInvalid
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return errEmailInvalid
	}

	return nil
}

func validatePassword(password string) error {
	if len(password) < passwordMinLength {
		return errPasswordTooShort
	}

	// @todo add checks for numbers, specials chars, etc...

	return nil
}
