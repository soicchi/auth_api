package controllers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRefreshTokenService struct {
	mock.Mock
}

func (m *MockRefreshTokenService) RefreshAccessToken(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

func TestPostRefreshToken(t *testing.T) {
	tests := []struct {
		name     string
		mock     func(mockRefreshTokenService *MockRefreshTokenService)
		wantBody string
		wantCode int
	}{
		{
			name: "success to refresh access token",
			mock: func(mockRefreshTokenService *MockRefreshTokenService) {
				mockRefreshTokenService.On("RefreshAccessToken", "refresh_token").Return("token", nil)
			},
			wantBody: "{\"data\":{\"access_token\":\"token\"},\"message\":\"New access token has been issued\"}\n",
			wantCode: http.StatusOK,
		},
		{
			name: "failed to refresh access token",
			mock: func(mockRefreshTokenService *MockRefreshTokenService) {
				mockRefreshTokenService.On("RefreshAccessToken", "refresh_token").Return("", fmt.Errorf("failed to verify refresh token"))
			},
			wantBody: "{\"data\":null,\"message\":\"bad request\"}\n",
			wantCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var mockRefreshTokenService MockRefreshTokenService
			test.mock(&mockRefreshTokenService)

			h := &RefreshTokenHandler{
				Service: &mockRefreshTokenService,
			}
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/key/refresh", nil)
			req.AddCookie(&http.Cookie{
				Name:  "refresh_token",
				Value: "refresh_token",
			})

			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			h.PostRefreshToken(ctx)
			assert.Equal(t, test.wantCode, rec.Code)
			assert.Equal(t, test.wantBody, rec.Body.String())
		})
	}
}
