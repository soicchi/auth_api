package middleware

import (
	"github.com/soicchi/auth_api/internal/utils"

	"github.com/labstack/echo/v4"
)

func JWTAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("access_token")
		if err != nil {
			return utils.UnauthorizedResponse(c, "invalid access token")
		}

		tokenString := cookie.Value
		if tokenString == "" {
			return utils.UnauthorizedResponse(c, "invalid access token")
		}

		if err := utils.ValidateJWT(tokenString); err != nil {
			return utils.UnauthorizedResponse(c, "invalid access token")
		}

		return next(c)
	}
}
