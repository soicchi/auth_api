package usecase

import (
	"testing"

	"github.com/soicchi/auth_api/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockDB struct {
	mock.Mock
}

type MockValidator struct {
	mock.Mock
}

func (m *MockDB) Create(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockValidator) Validate(value interface{}) error {
	args := m.Called(value)
	return args.Error(0)
}

func TestCreateUser(t *testing.T) {
	var mockDB MockDB
	var mockValidator MockValidator
	service := NewUserService(&mockDB, &mockValidator)
	email := "test@test.com"
	password := "password"
	user := models.NewUser(email, password)

	t.Run("create user successfully", func(t *testing.T) {
		mockValidator.On("Validate", user).Return(nil)
		mockDB.On("Create", user).Return(&gorm.DB{Error: nil})

		err := service.CreateUser(user.Email, user.Password)

		assert.NoError(t, err)
		mockValidator.AssertExpectations(t)
		mockDB.AssertExpectations(t)
	})

	// TODO: fix this test
	// t.Run("validation error", func(t *testing.T) {
	// 	mockValidator.On("Validate", user).Return(fmt.Errorf("validation error"))

	// 	err := service.CreateUser(email, password)

	// 	assert.Error(t, err)
	// 	assert.Contains(t, err.Error(), "validation error")
	// 	mockValidator.AssertExpectations(t)
	// })

	// t.Run("db error", func(t *testing.T) {
	// 	mockValidator.On("Validate", user).Return(nil)
	// 	mockDB.On("Create", user).Return(&gorm.DB{Error: fmt.Errorf("db error")})

	// 	err := service.CreateUser(user.Email, user.Password)

	// 	assert.Error(t, err)
	// 	assert.Contains(t, err.Error(), "db error")
	// 	mockValidator.AssertExpectations(t)
	// 	mockDB.AssertExpectations(t)
	// })
}

func TestNewUserService(t *testing.T) {
	var mockDB MockDB
	var mockValidator MockValidator
	service := NewUserService(&mockDB, &mockValidator)

	assert.IsType(t, &UserService{}, service)
	assert.Equal(t, &mockDB, service.DB)
	assert.Equal(t, &mockValidator, service.Validator)
}