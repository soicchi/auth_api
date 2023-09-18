package usecase

import (
	"fmt"

	"github.com/soicchi/auth_api/models"
	"github.com/soicchi/auth_api/utils"

	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
	Validator *utils.CustomValidator
}

func NewUserService(db *gorm.DB, cv *utils.CustomValidator) *UserService {
	return &UserService{
		DB:        db,
		Validator: cv,
	}
}

func (s *UserService) CreateUser(email string, password string) error {
	user := models.NewUser(email, password)

	if err := s.Validator.Validate(user); err != nil {
		return fmt.Errorf("error validating user struct %v", err)
	}

	if err := s.DB.Create(user).Error; err != nil {
		return fmt.Errorf("error creating user %v", err)
	}

	return nil
}
