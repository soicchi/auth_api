package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
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
