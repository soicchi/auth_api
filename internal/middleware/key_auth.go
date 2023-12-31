package middleware

import (
	"os"

	"github.com/soicchi/auth_api/internal/utils"

	"github.com/labstack/echo/v4"
)

func KeyAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		key := c.Request().Header.Get("API-KEY")
		if key == "" {
			return utils.UnauthorizedResponse(c, "Not found API-KEY value")
		}

		if key != os.Getenv("API_KEY") {
			return utils.UnauthorizedResponse(c, "Invalid API-KEY value")
		}

		return next(c)
	}
}
