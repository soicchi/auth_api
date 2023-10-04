package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTCustomClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken() (string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(tokenBytes), nil
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

func (c *JWTCustomClaims) GenerateJWT() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	return token.SignedString(jwtKey)
}

func ValidateJWT(tokenString string) error {
	token, err := parseJWT(tokenString)
	if err != nil {
		return fmt.Errorf("error parsing token %v", err)
	}

	_, ok := token.Claims.(*JWTCustomClaims)
	if !ok || !token.Valid {
		return fmt.Errorf("error validating token")
	}

	return nil
}

func parseJWT(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
}

func ExtractBearerToken(authHeader string) (string, error) {
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid Authorization header format")
	}

	return parts[1], nil
}
