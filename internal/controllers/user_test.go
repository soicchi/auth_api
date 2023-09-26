package controllers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/soicchi/auth_api/internal/utils"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func setup() *echo.Echo {
	e := echo.New()
	e.Validator = utils.NewCustomValidator()
	return e
}

func (m *MockUserService) CreateUser(email, password string) error {
	args := m.Called(email, password)
	return args.Error(0)
}

func TestNewUserHandler(t *testing.T) {
	service := &MockUserService{}
	handler := NewUserHandler(service)
	assert.NotNil(t, handler)
	assert.Equal(t, service, handler.Service)
}

func TestSignUpValid(t *testing.T) {
	var mockUserService MockUserService
	mockUserService.On("CreateUser", "test@test.com", "password").Return(nil)
	handler := &UserHandler{Service: &mockUserService}
	userJSON := `{"email": "test@test.com", "password": "password"}`

	e := setup()
	req := httptest.NewRequest(http.MethodPost, "/basic/signup", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	handler.SignUp(ctx)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "{\"data\":null,\"message\":\"Successfully created user\"}\n", rec.Body.String())
	mockUserService.AssertExpectations(t)
}

func TestSignWithBindError(t *testing.T) {
	var mockUserService MockUserService
	handler := &UserHandler{Service: &mockUserService}
	userJSON := `{"email": "test@test.com", "invalid": }`

	e := setup()
	req := httptest.NewRequest(http.MethodPost, "/basic/signup", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	handler.SignUp(ctx)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "{\"data\":null,\"message\":\"Invalid request\"}\n", rec.Body.String())
	mockUserService.AssertExpectations(t)
}

func TestSignUpWithEmailValidateError(t *testing.T) {
	var mockUserService MockUserService
	handler := &UserHandler{Service: &mockUserService}
	userJSON := `{"email": "test", "password": "password"}`

	e := setup()
	req := httptest.NewRequest(http.MethodPost, "/basic/signup", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	handler.SignUp(ctx)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "{\"data\":null,\"message\":\"Invalid request\"}\n", rec.Body.String())
	mockUserService.AssertExpectations(t)
}

func TestSignUpWithPasswordValidateError(t *testing.T) {
	var mockUserService MockUserService
	handler := &UserHandler{Service: &mockUserService}
	userJSON := `{"email": "test@test.com", "password": "pass"}`

	e := setup()
	req := httptest.NewRequest(http.MethodPost, "/basic/signup", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	handler.SignUp(ctx)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "{\"data\":null,\"message\":\"Invalid request\"}\n", rec.Body.String())
	mockUserService.AssertExpectations(t)
}

func TestSignUpWithCreateUserError(t *testing.T) {
	var mockUserService MockUserService
	mockUserService.On("CreateUser", "test@test.com", "password").Return(fmt.Errorf("error"))
	handler := &UserHandler{Service: &mockUserService}
	userJSON := `{"email": "test@test.com", "password": "password"}`

	e := setup()
	req := httptest.NewRequest(http.MethodPost, "/basic/signup", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	handler.SignUp(ctx)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, "{\"data\":null,\"message\":\"Failed to create user\"}\n", rec.Body.String())
	mockUserService.AssertExpectations(t)
}
