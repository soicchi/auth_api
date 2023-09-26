package usecase

import (
	"fmt"

	"github.com/soicchi/auth_api/internal/models"
)

type UserServiceImpl struct {
	Repo UserRepository
}

type UserRepository interface {
	CreateUser(user *models.User) error
}

func NewUserServiceImpl(repo UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		Repo: repo,
	}
}

func (s *UserServiceImpl) CreateUser(email string, password string) error {
	user := models.NewUser(email, password)

	err := s.Repo.CreateUser(user)
	if err != nil {
		return fmt.Errorf("error creating user %v", err)
	}

	return nil
}
