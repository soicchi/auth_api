package usecase

import (
	"fmt"
	"time"

	"github.com/soicchi/auth_api/internal/models"
	"github.com/soicchi/auth_api/internal/utils"
)

type RefreshTokenServiceImpl struct {
	TokenRepo RefreshTokenRepository
}

type RefreshTokenRepository interface {
	FetchByToken(token string) (models.RefreshToken, error)
}

func NewRefreshTokenServiceImpl(tokenRepo RefreshTokenRepository) *RefreshTokenServiceImpl {
	return &RefreshTokenServiceImpl{
		TokenRepo: tokenRepo,
	}
}

func (s *RefreshTokenServiceImpl) RefreshAccessToken(token string) (string, error) {
	refreshToken, err := verifyRefreshToken(s.TokenRepo, token)
	if err != nil {
		return "", fmt.Errorf("failed to verify refresh token: %w", err)
	}

	accessToken, err := utils.GenerateJWT(refreshToken.UserID)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	return accessToken, nil
}

func verifyRefreshToken(repo RefreshTokenRepository, token string) (models.RefreshToken, error) {
	refreshToken, err := repo.FetchByToken(token)
	if err != nil {
		return refreshToken, fmt.Errorf("failed to fetch refresh token: %w", err)
	}

	if refreshToken.Token == "" {
		return refreshToken, fmt.Errorf("refresh token not found")
	}

	if refreshToken.ExpiredAt.Before(time.Now()) {
		return refreshToken, fmt.Errorf("refresh token is expired")
	}

	return refreshToken, nil
}
