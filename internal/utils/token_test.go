package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "success generating token",
			wantErr: false,
		},
		{
			name:    "error generating token",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := GenerateToken()
			if test.wantErr && err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, got)
				assert.Equal(t, 64, len(got))
			}
		})
	}
}

func TestGenerateJWT(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	userID := uint(1)

	tokenString, err := GenerateJWT(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)
}

func TestExtractTokenFromHeader(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		wantErr bool
	}{
		{
			name:    "valid token",
			in:      "Bearer token",
			wantErr: false,
		},
		{
			name:    "invalid Header format",
			in:      "Bearerinvalid",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ExtractBearerToken(test.in)
			if test.wantErr && err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "token", got)
			}
		})
	}
}
