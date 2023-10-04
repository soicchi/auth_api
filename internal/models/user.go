package models

import (
	"errors"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string       `gorm:"unique;not null;size:255"`
	Password     string       `gorm:"not null;size:255"`
	RefreshToken RefreshToken `gorm:"constraint:OnDelete:CASCADE"`
}

type UserPostgresRepository struct {
	DB *gorm.DB
}

func NewUser(email, password string, token RefreshToken) *User {
	return &User{
		Email:        email,
		Password:     password,
		RefreshToken: token,
	}
}

func NewUserPostgresRepository(db *gorm.DB) *UserPostgresRepository {
	return &UserPostgresRepository{
		DB: db,
	}
}

func (r *UserPostgresRepository) CreateUser(user *User) (uint, error) {
	result := r.DB.Create(user)
	if result.Error != nil {
		return uint(0), result.Error
	}

	return user.ID, nil
}

func (r *UserPostgresRepository) FetchUserByEmail(email string) (*User, error) {
	var user User
	result := r.DB.Where("email = ?", email).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (r *UserPostgresRepository) FetchUsers() ([]User, error) {
	users := []User{}
	result := r.DB.Find(&users)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return users, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}
