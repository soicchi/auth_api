package usecase

import (
	"fmt"

	"github.com/soicchi/auth_api/internal/models"
	"github.com/soicchi/auth_api/internal/utils"
)

type UserServiceImpl struct {
	Repo UserRepository
}

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
}

func NewUserServiceImpl(repo UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		Repo: repo,
	}
}

func (s *UserServiceImpl) CreateUser(email string, password string) error {
	user := models.NewUser(email, password)

	// hash password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("error hashing password %v", err)
	}
	user.Password = hashedPassword

	if err := s.Repo.CreateUser(user); err != nil {
		return fmt.Errorf("error creating user %v", err)
	}

	return nil
}

func (s *UserServiceImpl) CheckSignIn(email, password string) error {
	user, err := s.Repo.GetUserByEmail(email)
	if err != nil {
		return fmt.Errorf("error getting user by email %v", err)
	}

	if !utils.ValidatePassword(user.Password, password) {
		return fmt.Errorf("error validating password")
	}

	return nil
}
