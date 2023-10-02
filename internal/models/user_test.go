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

func TestMain(m *testing.M) {
	// set up test database
	setupDB()

	// run tests
	code := m.Run()

	// tear down test database
	teardown()

	os.Exit(code)
}

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
	tests := []struct {
		name    string
		want    *User
		wantErr bool
	}{
		{
			name: "success creating user",
			want: &User{
				Email:    "test@test.com",
				Password: "password",
			},
			wantErr: false,
		},
		{
			name: "duplicate email error",
			want: &User{
				Email:    "test@test.com",
				Password: "password",
			},
			wantErr: true,
		},
	}

	// transaction
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := &UserPostgresRepository{
		DB: testDB,
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := repo.CreateUser(test.want)
			if test.wantErr && err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				var createdUser User
				testDB.First(&createdUser, test.want.ID)
				assert.Equal(t, test.want.Email, createdUser.Email)
				assert.Equal(t, test.want.Password, createdUser.Password)
			}
		})
	}
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
			wantErr: false,
		},
	}

	// transaction
	tx := testDB.Begin()
	defer tx.Rollback()

	// create test one user
	createTestUser()

	repo := &UserPostgresRepository{
		DB: testDB,
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := repo.GetUserByEmail(test.input)
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
		DB: testDB,
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := repo.GetUsers()
			if test.wantErr && err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.wantUsers, len(got))
			}
			if got != nil {
				testDB.Delete(got)
			}
		})
	}
}

// helper functions
// You have to build test-db using docker-compose
func setupDB() {
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
