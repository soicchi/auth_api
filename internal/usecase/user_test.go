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

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name string
		inputEmail string
		inputPassword string
		wantMock func(mockUserRepo *MockUserRepository)
		wantErr bool
	}{
		{
			name: "Valid create user",
			inputEmail: "test@test.com",
			inputPassword: "password",
			wantMock: func(mockUserRepo *MockUserRepository) {
				mockUserRepo.On("CreateUser", mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Create user with create error",
			inputEmail: "test@test.com",
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

			if err != nil || test.wantErr {
				assert.Error(t, err)
				assert.Equal(t, "error creating user db error", err.Error())
			} else {
				assert.NoError(t, err)
			}
			mockUserRepo.AssertExpectations(t)
		})
	}
}
