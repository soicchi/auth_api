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

func TestSignUp(t *testing.T) {
	tests := []struct {
		name string
		in string
		wantCode int
		wantBody string
		wantMock func(mockUserService *MockUserService)
	}{
		{
			name: "Valid signup",
			in: `{"email": "test@test.com", "password": "password"}`,
			wantCode: http.StatusOK,
			wantBody: "{\"data\":null,\"message\":\"Successfully created user\"}\n",
			wantMock: func(mockUserService *MockUserService) {
				mockUserService.On("CreateUser", "test@test.com", "password").Return(nil)
			},
		},
		{
			name: "Binding error",
			in: `{"email": "test@test.com", "invalid": }`,
			wantCode: http.StatusBadRequest,
			wantBody: "{\"data\":null,\"message\":\"Invalid request\"}\n",
			wantMock: func(mockUserService *MockUserService) {},
		},
		{
			name: "Email validation error",
			in: `{"email": "test", "password": "password"}`,
			wantCode: http.StatusBadRequest,
			wantBody: "{\"data\":null,\"message\":\"Invalid request\"}\n",
			wantMock: func(mockUserService *MockUserService) {},
		},
		{
			name: "Password validation error",
			in: `{"email": "test@test.com", "password": "pass"}`,
			wantCode: http.StatusBadRequest,
			wantBody: "{\"data\":null,\"message\":\"Invalid request\"}\n",
			wantMock: func(mockUserService *MockUserService) {},
		},
		{
			name: "Create user error",
			in: `{"email": "test@test.com", "password": "password"}`,
			wantCode: http.StatusInternalServerError,
			wantBody: "{\"data\":null,\"message\":\"Failed to create user\"}\n",
			wantMock: func(mockUserService *MockUserService) {
				mockUserService.On("CreateUser", "test@test.com", "password").Return(fmt.Errorf("error"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var mockUserService MockUserService
			test.wantMock(&mockUserService)
			handler := &UserHandler{Service: &mockUserService}

			e := echo.New()
			e.Validator = utils.NewCustomValidator()
			req := httptest.NewRequest(http.MethodPost, "/basic/signup", strings.NewReader(test.in))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			handler.SignUp(ctx)
			assert.Equal(t, test.wantCode, rec.Code)
			assert.Equal(t, test.wantBody, rec.Body.String())
			mockUserService.AssertExpectations(t)
		})
	}
}
