package usecase

import (
	"errors"
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

func TestCreateUserValid(t *testing.T) {
	var mockDB MockDB
	var mockValidator MockValidator
	user := models.NewUser("test@test.com", "password")
	mockDB.On("Create", user).Return(&gorm.DB{Error: nil})
	mockValidator.On("Validate", user).Return(nil)
	service := NewUserService(&mockDB, &mockValidator)

	t.Run("create user successfully", func(t *testing.T) {
		err := service.CreateUser(user.Email, user.Password)

		assert.NoError(t, err)
		mockValidator.AssertExpectations(t)
		mockDB.AssertExpectations(t)
	})
}

func TestCreateUserValidateError(t *testing.T) {
	var mockDB MockDB
	var mockValidator MockValidator
	user := models.NewUser("test", "password")
	mockValidator.On("Validate", user).Return(errors.New("validation error"))
	service := NewUserService(&mockDB, &mockValidator)

	t.Run("validation error", func(t *testing.T) {
		err := service.CreateUser(user.Email, user.Password)

		assert.Error(t, err)
		assert.Equal(t, "error validating user struct validation error", err.Error())
		mockValidator.AssertExpectations(t)
	})
}

func TestCreateUserDBError(t *testing.T) {
	var mockDB MockDB
	var mockValidator MockValidator

	user := models.NewUser("test@test.com", "password")
	mockValidator.On("Validate", user).Return(nil)
	mockDB.On("Create", user).Return(&gorm.DB{Error: errors.New("db error")})

	service := NewUserService(&mockDB, &mockValidator)

	t.Run("db error", func(t *testing.T) {
		err := service.CreateUser(user.Email, user.Password)

		assert.Error(t, err)
		assert.Equal(t, "error creating user db error", err.Error())
		mockValidator.AssertExpectations(t)
		mockDB.AssertExpectations(t)
	})
}

func TestNewUserService(t *testing.T) {
	var mockDB MockDB
	var mockValidator MockValidator
	service := NewUserService(&mockDB, &mockValidator)

	assert.IsType(t, &UserService{}, service)
	assert.Equal(t, &mockDB, service.DB)
	assert.Equal(t, &mockValidator, service.Validator)
}
