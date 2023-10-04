package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/soicchi/auth_api/internal/utils"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestJWTAuth(t *testing.T) {
	claims := &utils.JWTCustomClaims{UserID: 1}
	tokenString, _ := claims.GenerateJWT()
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
			req := httptest.NewRequest(http.MethodGet, "/jwt/users", nil)
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
