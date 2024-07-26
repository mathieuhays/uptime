package uptime

import (
	"errors"
	"testing"
)

func Test_validateEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  error
	}{
		{"valid email", "john@doe.com", nil},
		{"missing tld", "john@doe", errEmailInvalid},
		{"just a random string", "hey", errEmailInvalid},
		{"multi @", "john@doe@enterprise.net", errEmailInvalid},
		{"empty", "", errEmailInvalid},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateEmail(tt.email); !errors.Is(err, tt.want) {
				t.Errorf("validateEmail() error = %v, wantErr %v", err, tt.want)
			}
		})
	}
}

func Test_validatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     error
	}{
		{"valid password", "test123", nil},
		{"too short", "123", errPasswordTooShort},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validatePassword(tt.password); !errors.Is(err, tt.want) {
				t.Errorf("validatePassword() error = %v, wantErr %v", err, tt.want)
			}
		})
	}
}
