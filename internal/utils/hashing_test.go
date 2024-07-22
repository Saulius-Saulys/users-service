package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashAndSalt(t *testing.T) {
	tests := []struct {
		name             string
		originalPassword string
		comparePassword  string
		passwordsMatch   bool
	}{
		{
			name:             "original password and compare password matches",
			originalPassword: "password",
			comparePassword:  "password",
			passwordsMatch:   true,
		},
		{
			name:             "original password and compare password do not match",
			originalPassword: "password",
			comparePassword:  "password1",
			passwordsMatch:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashedPassword, err := HashAndSalt([]byte(tt.originalPassword))
			assert.NoError(t, err)
			assert.Equal(t, tt.passwordsMatch, comparePasswords(hashedPassword, []byte(tt.comparePassword)))
		})
	}
}
