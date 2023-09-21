package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	baseURI string = "/api/v1"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(email, password string) error {
	args := m.Called(email, password)
	return args.Error(0)
}

func TestNewUserController(t *testing.T) {
	uc := NewUserController(nil)
	assert.NotNil(t, uc)
	assert.IsType(t, &UserController{}, uc)
}

func TestSignUpValid(t *testing.T) {
	var userService MockUserService
	userService.On("CreateUser", "test@test.com", "password").Return(nil)
	uc := NewUserController(&userService)
	userJSON := `{"email": "test@test.com", "password": "password"}`

	t.Run("valid request", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, baseURI+"/basic/signup", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		uc.SignUp(ctx)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "{\"data\":null,\"message\":\"Successfully created user\"}\n", rec.Body.String())
	})
}

func TestSignUpCreateUserError(t *testing.T) {
	var userService MockUserService
	userService.On("CreateUser", "", "").Return(errors.New("error"))
	uc := NewUserController(&userService)
	userJSON := `{"email": "", "password": ""}`

	t.Run("failed create user", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, baseURI+"/basic/signup", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		uc.SignUp(ctx)
		assert.Equal(t, "{\"data\":null,\"message\":\"Failed to create user\"}\n", rec.Body.String())
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestSignBindError(t *testing.T) {
	var userService MockUserService
	uc := NewUserController(&userService)
	userJSON := `{"email": "test@test.com", "invalid": }`

	t.Run("bind error", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, baseURI+"/basic/signup", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		uc.SignUp(ctx)
		assert.Equal(t, "{\"data\":null,\"message\":\"Invalid request\"}\n", rec.Body.String())
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
