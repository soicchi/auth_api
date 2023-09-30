package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/soicchi/auth_api/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type JWTCustomClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func NewJWTCustomClaims(userID uint) *JWTCustomClaims {
	return &JWTCustomClaims{
		userID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth_api",
		},
	}
}

func (c *JWTCustomClaims) GenerateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	return token.SignedString(jwtKey)
}

func validateJWT(tokenString string) error {
	token, err := parseToken(tokenString)
	if err != nil {
		return fmt.Errorf("error parsing token %v", err)
	}

	_, ok := token.Claims.(*JWTCustomClaims)
	if !ok || !token.Valid {
		return fmt.Errorf("error validating token")
	}

	return nil
}

func parseToken(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
}

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

		if err := validateJWT(tokenString); err != nil {
			return utils.UnauthorizedResponse(c, "invalid access token")
		}

		return next(c)
	}
}
