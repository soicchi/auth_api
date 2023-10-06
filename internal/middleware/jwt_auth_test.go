package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/soicchi/auth_api/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestJWTAuth(t *testing.T) {
	userID := uint(1)
	tokenString, _ := utils.GenerateJWT(userID)

	tests := []struct {
		name       string
		in         string
		wantStatus int
	}{
		{
			name:       "valid token",
			in:         "Bearer " + tokenString,
			wantStatus: http.StatusOK,
		},
		{
			name:       "empty token",
			in:         "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "invalid token format",
			in:         "Bearerinvalidtoken",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "invalid token",
			in:         "Bearer invalid_token",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/jwt/users", nil)
			req.Header.Set("Authorization", test.in)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			middleware := JWTAuth(func(c echo.Context) error {
				return c.String(http.StatusOK, "test")
			})

			middleware(ctx)
			assert.Equal(t, test.wantStatus, rec.Code)
		})
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uint(1)
	tokenString, _ := utils.GenerateJWT(userID)
	tests := []struct {
		name    string
		in      string
		wantErr bool
	}{
		{
			name:    "valid token",
			in:      tokenString,
			wantErr: false,
		},
		{
			name:    "invalid token",
			in:      "invalid_token",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateJWT(test.in)
			if test.wantErr && err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParseJWT(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	userID := uint(1)
	tokenString, _ := utils.GenerateJWT(userID)
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
