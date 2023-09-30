package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

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

func TestGenerateToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")

	claims := &JWTCustomClaims{
		UserID: 1,
	}
	_, err := claims.GenerateToken()
	assert.NoError(t, err)
}

func TestValidateJWT(t *testing.T) {
	claims := &JWTCustomClaims{UserID: 1}
	tokenString, _ := claims.GenerateToken()
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
			err := validateJWT(test.jwt)
			if test.wantErr && err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParseToken(t *testing.T) {
	claims := &JWTCustomClaims{UserID: 1}
	tokenString, _ := claims.GenerateToken()
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
			_, err := parseToken(test.jwt)
			if test.wantErr && err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestJWTAuth(t *testing.T) {
	claims := &JWTCustomClaims{UserID: 1}
	tokenString, _ := claims.GenerateToken()
	tests := []struct {
		name       string
		cookie     *http.Cookie
		wantStatus int
	}{
		{
			name: "valid token",
			cookie: &http.Cookie{
				Name:  "access_token",
				Value: tokenString,
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "not Cookie",
			cookie:     nil,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "invalid token",
			cookie: &http.Cookie{
				Name:  "access_token",
				Value: "invalid_token",
			},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if test.cookie != nil {
				req.AddCookie(test.cookie)
			}

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
