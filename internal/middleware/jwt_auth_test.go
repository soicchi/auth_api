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
