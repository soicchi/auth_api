package models

import (
	"fmt"
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	// set up test database
	setupDB()
	defer teardown()

	// run tests
	code := m.Run()

	os.Exit(code)
}

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

	testDB.AutoMigrate(
		&User{},
		&RefreshToken{},
	)
}

func teardown() {
	testDB.Migrator().DropTable(
		&User{},
		&RefreshToken{},
	)
}
