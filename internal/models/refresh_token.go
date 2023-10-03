package models

import (
	"errors"
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

func NewRefreshTokenPostgresRepository(db *gorm.DB) RefreshTokenPostgresRepository {
	return RefreshTokenPostgresRepository{
		DB: db,
	}
}

func NewRefreshToken(userID uint, token string, expiredAt time.Time) RefreshToken {
	return RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiredAt: expiredAt,
	}
}

func (r *RefreshTokenPostgresRepository) CreateRefreshToken(refreshToken RefreshToken) error {
	result := r.DB.Create(&refreshToken)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *RefreshTokenPostgresRepository) FetchRefreshTokenByToken(token string) (RefreshToken, error) {
	var refreshToken RefreshToken
	result := r.DB.Where("token = ?", token).First(&refreshToken)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return refreshToken, nil
	}

	if result.Error != nil {
		return refreshToken, result.Error
	}

	return refreshToken, nil
}
