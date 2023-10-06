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

func GenerateToken() (string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(tokenBytes), nil
}

func GenerateJWT(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
		"sub":     "auth_api",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	return token.SignedString(jwtKey)
}

func ValidateJWT(tokenString string) error {
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

func ExtractBearerToken(authHeader string) (string, error) {
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid Authorization header format")
	}

	if parts[0] != "Bearer" {
		return "", fmt.Errorf("invalid Authorization header format")
	}

	return parts[1], nil
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
