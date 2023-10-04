package middleware

import (
	"log"

	"github.com/soicchi/auth_api/internal/utils"

	"github.com/labstack/echo/v4"
)

func JWTAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		authHeader := ctx.Request().Header.Get("Authorization")
		if authHeader == "" {
			log.Printf("Authorization header is empty")
			return utils.UnauthorizedResponse(ctx, "invalid access token")
		}

		tokenString, err := utils.ExtractBearerToken(authHeader)
		if err != nil {
			log.Printf("Failed to extract token from header: %v", err)
			return utils.UnauthorizedResponse(ctx, "invalid access token")
		}

		if err := utils.ValidateJWT(tokenString); err != nil {
			log.Printf("Failed to validate token: %v", err)
			return utils.UnauthorizedResponse(ctx, "invalid access token")
		}

		return next(ctx)
	}
}
