package middleware

import (
	"crypto/subtle"
	"os"

	"github.com/soicchi/auth_api/internal/utils"

	"github.com/labstack/echo/v4"
)

func BasicAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the Basic Authentication credentials from the request
		username, password, ok := c.Request().BasicAuth()

		if !ok {
			return utils.UnauthorizedResponse(c, "Not found Authorization header")
		}

		// Check credentials
		if !checkCredentials(username, password) {
			return utils.UnauthorizedResponse(c, "Invalid username or password")
		}

		return next(c)
	}
}

func checkCredentials(username, password string) bool {
	return subtle.ConstantTimeCompare([]byte(username), []byte(os.Getenv("BASIC_AUTH_USER"))) == 1 && subtle.ConstantTimeCompare([]byte(password), []byte(os.Getenv("BASIC_AUTH_PASSWORD"))) == 1
}
