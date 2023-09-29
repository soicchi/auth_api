package usecase

import (
	"fmt"
	"testing"

	"github.com/soicchi/auth_api/internal/models"
	"github.com/soicchi/auth_api/internal/utils"

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

func (m *MockUserRepository) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	return args.Get(0).(*models.User), args.Error(1)
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name          string
		inputEmail    string
		inputPassword string
		wantMock      func(mockUserRepo *MockUserRepository)
		wantErr       bool
	}{
		{
			name:          "Valid create user",
			inputEmail:    "test@test.com",
			inputPassword: "password",
			wantMock: func(mockUserRepo *MockUserRepository) {
				mockUserRepo.On("CreateUser", mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:          "Create user with create error",
			inputEmail:    "test@test.com",
			inputPassword: "password",
			wantMock: func(mockUserRepo *MockUserRepository) {
				mockUserRepo.On("CreateUser", mock.Anything).Return(fmt.Errorf("db error"))
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var mockUserRepo MockUserRepository
			test.wantMock(&mockUserRepo)
			userService := &UserServiceImpl{Repo: &mockUserRepo}
			user := &models.User{Email: test.inputEmail, Password: test.inputPassword}

			err := userService.CreateUser(user.Email, user.Password)

			if test.wantErr && err != nil {
				assert.Error(t, err)
				assert.Equal(t, "error creating user db error", err.Error())
			} else {
				assert.NoError(t, err)
			}
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestCheckSignIn(t *testing.T) {
	hashedPassword, _ := utils.HashPassword("password")

	tests := []struct {
		name          string
		inputEmail    string
		inputPassword string
		wantMock      func(mockUserRepo *MockUserRepository)
		ErrMsg        string
		wantErr       bool
	}{
		{
			name:          "Valid check sign in",
			inputEmail:    "test@test.com",
			inputPassword: "password",
			wantMock: func(mockUserRepo *MockUserRepository) {
				mockUserRepo.On("GetUserByEmail", "test@test.com").Return(&models.User{
					Email:    "test@test.com",
					Password: hashedPassword,
				}, nil)
			},
			ErrMsg:  "",
			wantErr: false,
		},
		{
			name:          "Check sign in with get user error",
			inputEmail:    "test@test.com",
			inputPassword: "password",
			wantMock: func(mockUserRepo *MockUserRepository) {
				mockUserRepo.On("GetUserByEmail", "test@test.com").Return(&models.User{}, fmt.Errorf("Not Found"))
			},
			ErrMsg:  "error getting user by email Not Found",
			wantErr: true,
		},
		{
			name:          "Check sign in with invalid password",
			inputEmail:    "test@test.com",
			inputPassword: "invalid",
			wantMock: func(mockUserRepo *MockUserRepository) {
				mockUserRepo.On("GetUserByEmail", "test@test.com").Return(&models.User{
					Email:    "test@test.com",
					Password: hashedPassword,
				}, nil)
			},
			ErrMsg:  "error validating password",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var mockUserRepo MockUserRepository
			test.wantMock(&mockUserRepo)
			userService := &UserServiceImpl{Repo: &mockUserRepo}

			err := userService.CheckSignIn(test.inputEmail, test.inputPassword)

			if test.wantErr && err != nil {
				assert.Error(t, err)
				assert.Equal(t, test.ErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
			mockUserRepo.AssertExpectations(t)
		})
	}
}
