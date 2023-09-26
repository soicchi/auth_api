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
	setENV()

	t.Run("valid username and password", func(t *testing.T) {
		ok := checkCredentials("test", "password")
		assert.True(t, ok)
	})

	t.Run("invalid username", func(t *testing.T) {
		ok := checkCredentials("invalid", "password")
		assert.False(t, ok)
	})

	t.Run("invalid password", func(t *testing.T) {
		ok := checkCredentials("test", "invalid")
		assert.False(t, ok)
	})
}

func TestBasicAuthValid(t *testing.T) {
	setENV()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/basic/signup", nil)
	req.Header.Set(echo.HeaderAuthorization, basicAuthHeader("test", "password"))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	middleware := BasicAuth(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	err := middleware(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestBasicAuthWithNotSetENV(t *testing.T) {
	os.Clearenv()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/basic/signup", nil)
	req.Header.Set(echo.HeaderAuthorization, basicAuthHeader("test", "password"))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	middleware := BasicAuth(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	err := middleware(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, "{\"data\":null,\"message\":\"Basic Auth credentials not set\"}\n", rec.Body.String())
}

func TestBasicAuthWithMissingAuthorizationHeader(t *testing.T) {
	setENV()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/basic/signup", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	middleware := BasicAuth(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	err := middleware(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Equal(t, "{\"data\":null,\"message\":\"Not found Authorization header\"}\n", rec.Body.String())
}

func TestBasicAuthWithInvalidUsername(t *testing.T) {
	setENV()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/basic/signup", nil)
	req.Header.Set(echo.HeaderAuthorization, basicAuthHeader("invalid", "password"))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	middleware := BasicAuth(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	err := middleware(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Equal(t, "{\"data\":null,\"message\":\"Invalid username or password\"}\n", rec.Body.String())
}

func TestBasicAuthWithInvalidPassword(t *testing.T) {
	setENV()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/basic/signup", nil)
	req.Header.Set(echo.HeaderAuthorization, basicAuthHeader("test", "invalid"))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	middleware := BasicAuth(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	err := middleware(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Equal(t, "{\"data\":null,\"message\":\"Invalid username or password\"}\n", rec.Body.String())
}

// Helper function
func basicAuthHeader(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

func setENV() {
	os.Setenv("BASIC_AUTH_USER", "test")
	os.Setenv("BASIC_AUTH_PASSWORD", "password")
}
