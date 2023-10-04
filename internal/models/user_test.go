package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNewUser(t *testing.T) {
	refreshToken := RefreshToken{
		UserID:    0,
		Token:     "token",
		ExpiredAt: time.Now().Add(time.Hour * 24 * 7),
	}
	user := NewUser("email", "password", refreshToken)
	assert.Equal(t, "email", user.Email)
	assert.Equal(t, "password", user.Password)
	assert.Equal(t, refreshToken, user.RefreshToken)
}

func TestNewUserRepository(t *testing.T) {
	var db *gorm.DB
	repo := NewUserPostgresRepository(db)
	assert.Equal(t, db, repo.DB)
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name         string
		in           *User
		wantEmail    string
		wantPassword string
		wantToken    string
		wantErr      bool
	}{
		{
			name: "success creating user",
			in: &User{
				Email:    "test@test.com",
				Password: "password",
				RefreshToken: RefreshToken{
					UserID:    uint(0),
					Token:     "token",
					ExpiredAt: time.Now().Add(time.Hour * 24 * 7),
				},
			},
			wantEmail:    "test@test.com",
			wantPassword: "password",
			wantToken:    "token",
			wantErr:      false,
		},
		{
			name: "duplicate email error",
			in: &User{
				Email:    "test@test.com",
				Password: "password",
				RefreshToken: RefreshToken{
					UserID:    0,
					Token:     "token",
					ExpiredAt: time.Now().Add(time.Hour * 24 * 7),
				},
			},
			wantEmail:    "",
			wantPassword: "",
			wantToken:    "",
			wantErr:      true,
		},
	}

	// transaction
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := &UserPostgresRepository{
		DB: tx,
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userID, err := repo.CreateUser(test.in)
			if test.wantErr && err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				var createdUser User
				tx.Preload("RefreshToken").Where("email = ?", test.wantEmail).First(&createdUser)
				assert.Equal(t, test.wantEmail, createdUser.Email)
				assert.Equal(t, test.wantPassword, createdUser.Password)
				assert.Equal(t, test.wantToken, createdUser.RefreshToken.Token)
				assert.NotEmpty(t, userID)
			}
		})
	}
}

func TestFetUserByEmail(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *User
		wantErr bool
	}{
		{
			name:  "success",
			input: "test@test.com",
			want: &User{
				Email:    "test@test.com",
				Password: "password",
			},
			wantErr: false,
		},
		{
			name:    "not found",
			input:   "invalid@test.com",
			want:    nil,
			wantErr: false,
		},
	}

	// transaction
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := &UserPostgresRepository{
		DB: tx,
	}

	// create user
	user := &User{
		Email:    "test@test.com",
		Password: "password",
		RefreshToken: RefreshToken{
			UserID:    0,
			Token:     "token",
			ExpiredAt: time.Now().Add(time.Hour * 24 * 7),
		},
	}
	repo.DB.Create(user)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := repo.FetchUserByEmail(test.input)
			if test.wantErr && err != nil {
				assert.Error(t, err)
			} else if test.want != nil {
				assert.NoError(t, err)
				assert.Equal(t, test.want.Email, got.Email)
				assert.Equal(t, test.want.Password, got.Password)
			} else {
				assert.NoError(t, err)
				assert.Nil(t, got)
			}
		})
	}
}

func TestGetUsers(t *testing.T) {
	tests := []struct {
		name      string
		wantUsers int
		wantErr   bool
	}{
		{
			name:      "success getting users",
			wantUsers: 1,
			wantErr:   false,
		},
		{
			name:      "no users",
			wantUsers: 0,
			wantErr:   false,
		},
		{
			name:      "error getting users",
			wantUsers: 0,
			wantErr:   true,
		},
	}

	// transaction
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := &UserPostgresRepository{
		DB: tx,
	}

	// create user
	user := &User{
		Email:    "test@test.com",
		Password: "password",
		RefreshToken: RefreshToken{
			UserID:    0,
			Token:     "token",
			ExpiredAt: time.Now().Add(time.Hour * 24 * 7),
		},
	}
	repo.DB.Create(user)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := repo.FetchUsers()
			if test.wantErr && err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.wantUsers, len(got))
			}
			if got != nil {
				tx.Delete(got)
			}
		})
	}
}
