package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique;not null;size:255"`
	Password string `gorm:"not null;size:255"`
}

type UserPostgresRepository struct {
	DB *gorm.DB
}

func NewUser(email, password string) *User {
	return &User{
		Email:    email,
		Password: password,
	}
}

func NewUserPostgresRepository(db *gorm.DB) *UserPostgresRepository {
	return &UserPostgresRepository{
		DB: db,
	}
}

func (r *UserPostgresRepository) CreateUser(user *User) error {
	result := r.DB.Create(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *UserPostgresRepository) GetUserByEmail(email string) (*User, error) {
	var user User
	result := r.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
