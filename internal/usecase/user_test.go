package usecase

import (
	"errors"
	"testing"

	"github.com/soicchi/auth_api/internal/models"
	"github.com/soicchi/auth_api/internal/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"github.com/go-playground/validator/v10"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Create(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func TestCreateUserValid(t *testing.T) {
	var mockDB MockDB
	cv := utils.CustomValidator{Validator: validator.New()}
	user := &models.User{Email: "test@test.com", Password: "password"}
	mockDB.On("Create", user).Return(&gorm.DB{Error: nil})
	service := &UserService{
		DB:        &mockDB,
		Validator: &cv,
	}

	t.Run("create user successfully", func(t *testing.T) {
		err := service.CreateUser(user.Email, user.Password)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})
}

func TestCreateUserValidateError(t *testing.T) {
	var mockDB MockDB
	cv := utils.CustomValidator{Validator: validator.New()}
	user := &models.User{Email: "test", Password: "password"}
	service := &UserService{
		DB:        &mockDB,
		Validator: &cv,
	}

	t.Run("validation error", func(t *testing.T) {
		err := service.CreateUser(user.Email, user.Password)

		assert.Error(t, err)
		assert.Equal(t, "error validating user struct error validating struct Key: 'User.Email' Error:Field validation for 'Email' failed on the 'email' tag", err.Error())
	})
}

func TestCreateUserDBError(t *testing.T) {
	var mockDB MockDB
	cv := utils.CustomValidator{Validator: validator.New()}

	user := &models.User{Email: "test@test.com", Password: "password"}
	mockDB.On("Create", user).Return(&gorm.DB{Error: errors.New("db error")})
	service := &UserService{
		DB:        &mockDB,
		Validator: &cv,
	}

	t.Run("db error", func(t *testing.T) {
		err := service.CreateUser(user.Email, user.Password)

		assert.Error(t, err)
		assert.Equal(t, "error creating user db error", err.Error())
		mockDB.AssertExpectations(t)
	})
}

func TestNewUserService(t *testing.T) {
	var mockDB MockDB
	cv := utils.CustomValidator{Validator: validator.New()}
	service := NewUserService(&mockDB, &cv)

	assert.IsType(t, &UserService{}, service)
	assert.Equal(t, &mockDB, service.DB)
	assert.Equal(t, &cv, service.Validator)
}
