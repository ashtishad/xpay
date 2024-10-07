package secure

import (
	"errors"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestGeneratePasswordHash(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  error
	}{
		{
			name:     "Empty password",
			password: "",
			wantErr:  errEmptyPassword,
		},
		{
			name:     "Password too short",
			password: "short",
			wantErr:  errPasswordTooShort,
		},
		{
			name:     "Password at minimum length",
			password: "12345678",
			wantErr:  nil,
		},
		{
			name:     "Password exceeds maximum length for bcrypt",
			password: strings.Repeat("a", 73),
			wantErr:  errPasswordTooLong,
		},
		{
			name:     "Valid Password with special characters",
			password: "P@ssw0rd!@#$%^&*()",
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GeneratePasswordHash(tt.password)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GeneratePasswordHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil {
				if got == "" {
					t.Error("GeneratePasswordHash() returned empty string for valid password")
					return
				}

				if err := bcrypt.CompareHashAndPassword([]byte(got), []byte(tt.password)); err != nil {
					t.Errorf("Generated hash does not match original password: %v", err)
				}
			}
		})
	}
}
