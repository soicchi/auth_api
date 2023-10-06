package utils

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func TestValidateJWT(t *testing.T) {
	userID := uint(1)
	tokenString, _ := GenerateJWT(userID)
	tests := []struct {
		name    string
		jwt     string
		wantErr bool
	}{
		{
			name:    "valid token",
			jwt:     tokenString,
			wantErr: false,
		},
		{
			name:    "invalid token",
			jwt:     "invalid_token",
			wantErr: true,
		},
	}

	os.Setenv("JWT_SECRET", "test_secret")

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateJWT(test.jwt)
			if test.wantErr && err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParseJWT(t *testing.T) {
	userID := uint(1)
	tokenString, _ := GenerateJWT(userID)
	tests := []struct {
		name    string
		jwt     string
		wantErr bool
	}{
		{
			name:    "valid token",
			jwt:     tokenString,
			wantErr: false,
		},
		{
			name:    "invalid token",
			jwt:     "invalid_token",
			wantErr: true,
		},
	}

	os.Setenv("JWT_SECRET", "test_secret")

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := parseJWT(test.jwt)
			if test.wantErr && err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
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

func TestCheckTokenExpiration(t *testing.T) {
	tests := []struct {
		name    string
		in      jwt.MapClaims
		wantErr bool
	}{
		{
			name: "valid token",
			in: jwt.MapClaims{
				"user_id": uint(1),
				"exp":     float64(time.Now().Add(time.Hour * 1).Unix()),
				"sup":     "auth_api",
			},
			wantErr: false,
		},
		{
			name:    "not include exp token",
			in:      jwt.MapClaims{},
			wantErr: true,
		},
		{
			name: "expired token",
			in: jwt.MapClaims{
				"user_id": uint(1),
				"exp":     float64(time.Now().Add(time.Hour * -1).Unix()),
				"sup":     "auth_api",
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := checkTokenExpiration(test.in)
			if test.wantErr && err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
