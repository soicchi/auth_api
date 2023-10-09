package models

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model
	UserID    uint      `gorm:"not null"`
	Token     string    `gorm:"not null"`
	ExpiredAt time.Time `gorm:"not null"`
}

type RefreshTokenPostgresRepository struct {
	DB *gorm.DB
}

func NewRefreshTokenPostgresRepository(db *gorm.DB) *RefreshTokenPostgresRepository {
	return &RefreshTokenPostgresRepository{
		DB: db,
	}
}

func NewRefreshToken(token string) RefreshToken {
	return RefreshToken{
		Token:     token,
		ExpiredAt: time.Now().Add(time.Hour * 24 * 7),
	}
}

func (r *RefreshTokenPostgresRepository) FetchByToken(token string) (RefreshToken, error) {
	var refreshToken RefreshToken
	result := r.DB.Where("token = ?", token).First(&refreshToken)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return refreshToken, nil
	}

	if result.Error != nil {
		return refreshToken, fmt.Errorf("failed to fetch refresh token: %w", result.Error)
	}

	return refreshToken, nil
}
