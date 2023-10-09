package usecase

import (
	"fmt"
	"testing"

	"github.com/soicchi/auth_api/internal/models"
	"github.com/soicchi/auth_api/internal/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *models.User) (uint, error) {
	args := m.Called(user)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockUserRepository) FetchUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FetchUsers() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
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
				mockUserRepo.On("CreateUser", mock.Anything).Return(uint(1), nil)
			},
			wantErr: false,
		},
		{
			name:          "Create user with create error",
			inputEmail:    "test@test.com",
			inputPassword: "password",
			wantMock: func(mockUserRepo *MockUserRepository) {
				mockUserRepo.On("CreateUser", mock.Anything).Return(uint(0), fmt.Errorf("db error"))
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var mockUserRepo MockUserRepository
			var mockTokenRepo MockRefreshTokenRepository
			test.wantMock(&mockUserRepo)
			userService := &UserServiceImpl{
				UserRepo:  &mockUserRepo,
				TokenRepo: &mockTokenRepo,
			}

			tokens, err := userService.CreateUser(test.inputEmail, test.inputPassword)

			if test.wantErr && err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, tokens)
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
				mockUserRepo.On("FetchUserByEmail", "test@test.com").Return(&models.User{
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
				mockUserRepo.On("FetchUserByEmail", "test@test.com").Return(&models.User{}, fmt.Errorf("Not Found"))
			},
			ErrMsg:  "Not Found",
			wantErr: true,
		},
		{
			name:          "Check sign in with invalid password",
			inputEmail:    "test@test.com",
			inputPassword: "invalid",
			wantMock: func(mockUserRepo *MockUserRepository) {
				mockUserRepo.On("FetchUserByEmail", "test@test.com").Return(&models.User{
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
			var mockTokenRepo MockRefreshTokenRepository
			test.wantMock(&mockUserRepo)
			userService := &UserServiceImpl{
				UserRepo:  &mockUserRepo,
				TokenRepo: &mockTokenRepo,
			}

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

func TestFetchAllUsers(t *testing.T) {
	tests := []struct {
		name     string
		wantMock func(mockUserRepo *MockUserRepository)
		ErrMsg   string
		wantErr  bool
	}{
		{
			name: "Valid fetch all users",
			wantMock: func(mockUserRepo *MockUserRepository) {
				mockUserRepo.On("FetchUsers").Return([]models.User{
					{
						Model: gorm.Model{
							ID: 1,
						},
						Email: "test@test.com",
					},
					{
						Model: gorm.Model{
							ID: 2,
						},
						Email: "test2@test.com",
					},
				}, nil)
			},
			ErrMsg:  "",
			wantErr: false,
		},
		{
			name: "Fetch all users with get users error",
			wantMock: func(mockUserRepo *MockUserRepository) {
				mockUserRepo.On("FetchUsers").Return([]models.User{}, fmt.Errorf("Not Found"))
			},
			ErrMsg:  "Not Found",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var mockUserRepo MockUserRepository
			var mockTokenRepo MockRefreshTokenRepository
			test.wantMock(&mockUserRepo)
			userService := &UserServiceImpl{
				UserRepo:  &mockUserRepo,
				TokenRepo: &mockTokenRepo,
			}

			users, err := userService.FetchAllUsers()

			if test.wantErr && err != nil {
				assert.Error(t, err)
				assert.Equal(t, test.ErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, users)
			}
			mockUserRepo.AssertExpectations(t)
		})
	}
}
