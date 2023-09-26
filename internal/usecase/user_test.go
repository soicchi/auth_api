package usecase

import (
	"fmt"
	"testing"

	"github.com/soicchi/auth_api/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestCreateUserValid(t *testing.T) {
	var mockUserRepo MockUserRepository
	user := &models.User{Email: "test@test.com", Password: "password"}
	mockUserRepo.On("CreateUser", mock.Anything).Return(nil)
	userService := &UserServiceImpl{Repo: &mockUserRepo}

	err := userService.CreateUser(user.Email, user.Password)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

func TestCreateUserWithCreateError(t *testing.T) {
	var mockUserRepo MockUserRepository
	user := &models.User{Email: "test@test.com", Password: "password"}
	mockUserRepo.On("CreateUser", mock.Anything).Return(fmt.Errorf("db error"))
	userService := &UserServiceImpl{Repo: &mockUserRepo}

	err := userService.CreateUser(user.Email, user.Password)

	assert.Error(t, err)
	assert.Equal(t, "error creating user db error", err.Error())
	mockUserRepo.AssertExpectations(t)
}
