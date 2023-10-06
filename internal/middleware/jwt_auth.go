package middleware

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/soicchi/auth_api/internal/utils"

	"github.com/golang-jwt/jwt/v5"
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

		if err := validateJWT(tokenString); err != nil {
			log.Printf("Failed to validate token: %v", err)
			return utils.UnauthorizedResponse(ctx, "invalid access token")
		}

		return next(ctx)
	}
}

func validateJWT(tokenString string) error {
	token, err := parseJWT(tokenString)
	if err != nil {
		return fmt.Errorf("error parsing token %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return fmt.Errorf("invalid token")
	}

	// check token expiration
	if err := checkTokenExpiration(claims); err != nil {
		return fmt.Errorf("error checking token expiration %v", err)
	}

	return nil
}

func parseJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
}

func checkTokenExpiration(claims jwt.MapClaims) error {
	exp, ok := claims["exp"].(float64)
	if !ok {
		return fmt.Errorf("error getting token expiration")
	}

	if time.Now().Unix() > int64(exp) {
		return fmt.Errorf("token expired")
	}

	return nil
}
