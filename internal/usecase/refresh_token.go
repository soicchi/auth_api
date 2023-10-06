package usecase

import (
	"fmt"
	"time"

	"github.com/soicchi/auth_api/internal/models"
)

type RefreshTokenServiceImpl struct {
	TokenRepo RefreshTokenRepository
}

type RefreshTokenRepository interface {
	FetchByToken(token string) (models.RefreshToken, error)
}

func (s *RefreshTokenServiceImpl) VerifyRefreshToken(token string) error {
	refreshToken, err := s.TokenRepo.FetchByToken(token)
	if err != nil {
		return fmt.Errorf("failed to fetch refresh token: %w", err)
	}

	if refreshToken.Token == "" {
		return fmt.Errorf("refresh token not found")
	}

	if refreshToken.ExpiredAt.Before(time.Now()) {
		return fmt.Errorf("refresh token is expired")
	}

	return nil
}
