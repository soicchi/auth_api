package usecase

import (
	"fmt"

	"github.com/soicchi/auth_api/models"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB, validator *validator.Validate) *UserService {
	return &UserService{
		DB:        db,
		Validator: validator,
	}
}

func (s *UserService) CreateUser(email string, password string) error {
	user := models.NewUser(email, password)

	if err := s.Validator.Struct(user); err != nil {
		return fmt.Errorf("error validating user struct %v", err)
	}

	if err := s.DB.Create(user).Error; err != nil {
		return fmt.Errorf("error creating user %v", err)
	}

	return nil
}
