package models

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func TestNewUser(t *testing.T) {
	user := NewUser("email", "password")
	assert.Equal(t, "email", user.Email)
	assert.Equal(t, "password", user.Password)
}

func TestNewUserRepository(t *testing.T) {
	var db *gorm.DB
	repo := NewUserPostgresRepository(db)
	assert.Equal(t, db, repo.DB)
}

func TestCreateUser(t *testing.T) {
	setup()
	defer teardown()

	user := &User{
		Email:    "test@test.com",
		Password: "password",
	}
	repo := &UserPostgresRepository{
		DB: testDB,
	}
	err := repo.CreateUser(user)
	assert.NoError(t, err)

	// Check if user is created
	var createdUser User
	testDB.First(&createdUser, user.ID)
	assert.Equal(t, user.Email, createdUser.Email)
	assert.Equal(t, user.Password, createdUser.Password)

	// duplicate email error
	dupUser := &User{
		Email:    "test@test.com",
		Password: "password",
	}
	err = repo.CreateUser(dupUser)
	assert.Error(t, err)
}

func TestGetUserByEmail(t *testing.T) {
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
			wantErr: true,
		},
	}

	setup()
	defer teardown()

	createTestUser()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &UserPostgresRepository{
				DB: testDB,
			}
			got, err := repo.GetUserByEmail(test.input)
			if test.wantErr && err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want.Email, got.Email)
				assert.Equal(t, test.want.Password, got.Password)
			}
		})
	}
}

// helper functions
// You have to build test-db using docker-compose
func setup() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("TEST_DB_HOST"), os.Getenv("TEST_DB_USER"), os.Getenv("TEST_DB_PASSWORD"), os.Getenv("TEST_DB_NAME"), os.Getenv("DB_PORT"), os.Getenv("DB_SSLMODE"),
	)
	var err error
	testDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	testDB.AutoMigrate(&User{})
}

func teardown() {
	testDB.Migrator().DropTable(&User{})
}

func createTestUser() {
	user := &User{
		Email:    "test@test.com",
		Password: "password",
	}
	repo := &UserPostgresRepository{DB: testDB}
	repo.CreateUser(user)
}
