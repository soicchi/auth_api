package usecase

import (
	"fmt"

	"github.com/soicchi/auth_api/internal/models"

	"gorm.io/gorm"
)

type IFDatabase interface {
	Create(value interface{}) *gorm.DB
}

type IFValidator interface {
	Validate(value interface{}) error
}

type UserService struct {
	DB        IFDatabase
	Validator IFValidator
}

func NewUserService(db IFDatabase, v IFValidator) *UserService {
	return &UserService{
		DB:        db,
		Validator: v,
	}
}

func (s *UserService) CreateUser(email string, password string) error {
	user := models.NewUser(email, password)

	if err := s.Validator.Validate(user); err != nil {
		return fmt.Errorf("error validating user struct %v", err)
	}

	result := s.DB.Create(user)
	if result.Error != nil {
		return fmt.Errorf("error creating user %v", result.Error)
	}

	return nil
}
