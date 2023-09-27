package middleware

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCheckBasicAuth(t *testing.T) {
	tests := []struct {
		name string
		inputUsername string
		inputPassword string
		want bool
	}{
		{
			name: "Valid username and password",
			inputUsername: "test",
			inputPassword: "password",
			want: true,
		},
		{
			name: "Invalid username",
			inputUsername: "invalid",
			inputPassword: "password",
			want: false,
		},
		{
			name: "Invalid password",
			inputUsername: "test",
			inputPassword: "invalid",
			want: false,
		},
	}

	// Set environment variables
	os.Setenv("BASIC_AUTH_USER", "test")
	os.Setenv("BASIC_AUTH_PASSWORD", "password")

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ok := checkCredentials(test.inputUsername, test.inputPassword)
			assert.Equal(t, test.want, ok)
		})
	}
}

func TestBasicAuth(t *testing.T) {
	tests := []struct {
		name string
		inputUsername string
		inputPassword string
		wantCode int
		wantBody string
		isSetENV bool
		isSetHeader bool
	}{
		{
			name: "Valid basic auth",
			inputUsername: "test",
			inputPassword: "password",
			wantCode: http.StatusOK,
			wantBody: "test",
			isSetENV: true,
			isSetHeader: true,
		},
		{
			name: "Not set environment variables",
			inputUsername: "test",
			inputPassword: "password",
			wantCode: http.StatusInternalServerError,
			wantBody: "{\"data\":null,\"message\":\"Basic Auth credentials not set\"}\n",
			isSetENV: false,
			isSetHeader: true,
		},
		{
			name: "Missing Authorization header",
			inputUsername: "test",
			inputPassword: "password",
			wantCode: http.StatusUnauthorized,
			wantBody: "{\"data\":null,\"message\":\"Not found Authorization header\"}\n",
			isSetENV: true,
			isSetHeader: false,
		},
		{
			name: "Invalid username",
			inputUsername: "invalid",
			inputPassword: "password",
			wantCode: http.StatusUnauthorized,
			wantBody: "{\"data\":null,\"message\":\"Invalid username or password\"}\n",
			isSetENV: true,
			isSetHeader: true,
		},
		{
			name: "Invalid password",
			inputUsername: "test",
			inputPassword: "invalid",
			wantCode: http.StatusUnauthorized,
			wantBody: "{\"data\":null,\"message\":\"Invalid username or password\"}\n",
			isSetENV: true,
			isSetHeader: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.isSetENV {
				os.Setenv("BASIC_AUTH_USER", "test")
				os.Setenv("BASIC_AUTH_PASSWORD", "password")
			} else {
				os.Clearenv()
			}

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/basic/signup", nil)
			if test.isSetHeader {
				req.Header.Set(echo.HeaderAuthorization, basicAuthHeader(test.inputUsername, test.inputPassword))
			}

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Assume Basic authentication returns status 200 on success
			middleware := BasicAuth(func(c echo.Context) error {
				return c.String(http.StatusOK, "test")
			})

			middleware(c)
			assert.Equal(t, test.wantCode, rec.Code)
			assert.Equal(t, test.wantBody, rec.Body.String())
		})
	}
}

// Helper function
func basicAuthHeader(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}
