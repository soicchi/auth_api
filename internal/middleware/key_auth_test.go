package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestKeyAuth(t *testing.T) {
	tests := []struct {
		name     string
		inputKey string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid key auth",
			inputKey: "testkey",
			wantCode: http.StatusOK,
			wantBody: "test",
		},
		{
			name:     "Invalid key auth",
			inputKey: "invalid",
			wantCode: http.StatusUnauthorized,
			wantBody: "{\"data\":null,\"message\":\"Invalid API-KEY value\"}\n",
		},
	}

	os.Setenv("API_KEY", "testkey")

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("API-KEY", test.inputKey)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Assume Basic authentication returns status 200 on success
			middleware := KeyAuth(func(c echo.Context) error {
				return c.String(http.StatusOK, "test")
			})

			middleware(c)
			assert.Equal(t, test.wantCode, rec.Code)
			assert.Equal(t, test.wantBody, rec.Body.String())
		})
	}
}
