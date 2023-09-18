package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
}

func NewUser(email string, password string) *User {
	return &User{
		Email:    email,
		Password: password,
	}
}
