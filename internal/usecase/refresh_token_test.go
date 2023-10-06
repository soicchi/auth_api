package usecase

import (
	"fmt"
	"testing"
	"time"

	"github.com/soicchi/auth_api/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRefreshTokenRepository struct {
	mock.Mock
}

func (m *MockRefreshTokenRepository) FetchByToken(token string) (models.RefreshToken, error) {
	args := m.Called(token)
	return args.Get(0).(models.RefreshToken), args.Error(1)
}

func TestVerifyRefreshToken(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(mockTokenRepo *MockRefreshTokenRepository)
		wantErr bool
	}{
		{
			name: "success",
			mock: func(mockTokenRepo *MockRefreshTokenRepository) {
				mockTokenRepo.On("FetchByToken", "token").Return(models.RefreshToken{
					UserID:    1,
					Token:     "token",
					ExpiredAt: time.Now().Add(time.Hour * 1),
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "failed to fetch refresh token",
			mock: func(mockTokenRepo *MockRefreshTokenRepository) {
				mockTokenRepo.On("FetchByToken", "token").Return(models.RefreshToken{}, fmt.Errorf("failed to fetch refresh token"))
			},
			wantErr: true,
		},
		{
			name: "refresh token not found",
			mock: func(mockTokenRepo *MockRefreshTokenRepository) {
				mockTokenRepo.On("FetchByToken", "token").Return(models.RefreshToken{}, nil)
			},
			wantErr: true,
		},
		{
			name: "refresh token is expired",
			mock: func(mockTokenRepo *MockRefreshTokenRepository) {
				mockTokenRepo.On("FetchByToken", "token").Return(models.RefreshToken{
					UserID:    1,
					Token:     "token",
					ExpiredAt: time.Now().Add(time.Hour * -1),
				}, nil)
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var mockTokenRepo MockRefreshTokenRepository
			test.mock(&mockTokenRepo)
			tokenService := &RefreshTokenServiceImpl{
				TokenRepo: &mockTokenRepo,
			}
			err := tokenService.VerifyRefreshToken("token")
			if test.wantErr && err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
