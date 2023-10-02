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
	FetchUserByEmail(email string) (*models.User, error)
	FetchUsers() ([]models.User, error)
}

type ResponseUser struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

type AllUsersResponse struct {
	Users []ResponseUser `json:"users"`
}

func NewUserServiceImpl(repo UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		Repo: repo,
	}
}

func newResponseUser(user *models.User) ResponseUser {
	return ResponseUser{
		ID:    user.ID,
		Email: user.Email,
	}
}

func newAllUsersResponse(users []models.User) AllUsersResponse {
	responseUsers := make([]ResponseUser, len(users))
	for _, user := range users {
		responseUser := newResponseUser(&user)
		responseUsers = append(responseUsers, responseUser)
	}

	return AllUsersResponse{
		Users: responseUsers,
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
	user, err := s.Repo.FetchUserByEmail(email)
	if err != nil {
		return fmt.Errorf("error getting user by email %v", err)
	}

	if !utils.ValidatePassword(user.Password, password) {
		return fmt.Errorf("error validating password")
	}

	return nil
}

func (s *UserServiceImpl) FetchAllUsers() (AllUsersResponse, error) {
	allUsersResponse := AllUsersResponse{}
	users, err := s.Repo.FetchUsers()
	if err != nil {
		return allUsersResponse, fmt.Errorf("error getting users %v", err)
	}

	allUsersResponse = newAllUsersResponse(users)

	return allUsersResponse, nil
}
