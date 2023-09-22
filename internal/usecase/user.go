package usecase

import (
	"fmt"

	"github.com/soicchi/auth_api/internal/models"
	"github.com/soicchi/auth_api/internal/utils"

	"gorm.io/gorm"
)

type IFDatabase interface {
	Create(value interface{}) *gorm.DB
}

type UserService struct {
	DB        IFDatabase
	Validator *utils.CustomValidator 
}

func NewUserService(db IFDatabase, validator *utils.CustomValidator) *UserService {
	return &UserService{
		DB:        db,
		Validator: validator,
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
