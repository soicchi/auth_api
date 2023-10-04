package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNewRefreshToken(t *testing.T) {
	refreshToken := NewRefreshToken("token")
	assert.Equal(t, uint(0), refreshToken.UserID)
	assert.Equal(t, "token", refreshToken.Token)
}

func TestNewRefreshTokenRepository(t *testing.T) {
	var db *gorm.DB
	repo := NewRefreshTokenPostgresRepository(db)
	assert.Equal(t, db, repo.DB)
}

func TestFetchRefreshTokenByToken(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name: "success fetching refresh token",
			want: "token",
			wantErr: false,
		},
		{
			name:    "refresh token not found",
			want:    "",
			wantErr: false,
		},
	}

	// transaction
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := RefreshTokenPostgresRepository{
		DB: tx,
	}

	// create user
	user := &User{
		Email: "test@test.com",
		Password: "password",
		RefreshToken: RefreshToken{
			UserID:    0,
			Token:	"token",
			ExpiredAt: time.Now().Add(time.Hour * 24 * 7),
		},
	}
	repo.DB.Create(user)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := repo.FetchByToken(test.want)
			if test.wantErr && err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got.Token)
			}
		})
	}
}
