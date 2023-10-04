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

func TestNewJWTCustomClaims(t *testing.T) {
	var testUserID uint = 1
	customClaims := NewJWTCustomClaims(testUserID)

	assert.Equal(t, testUserID, customClaims.UserID)

	expectedExp := jwt.NewNumericDate(time.Now().Add(1 * time.Hour))
	assert.Equal(t, expectedExp, customClaims.ExpiresAt)

	expectedIssuedAt := jwt.NewNumericDate(time.Now())
	assert.Equal(t, expectedIssuedAt, customClaims.IssuedAt)
	assert.Equal(t, "auth_api", customClaims.Issuer)
}

func TestGenerateJWT(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")

	claims := &JWTCustomClaims{
		UserID: 1,
	}
	_, err := claims.GenerateJWT()
	assert.NoError(t, err)
}

func TestValidateJWT(t *testing.T) {
	claims := &JWTCustomClaims{UserID: 1}
	tokenString, _ := claims.GenerateJWT()
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
	claims := &JWTCustomClaims{UserID: 1}
	tokenString, _ := claims.GenerateJWT()
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
