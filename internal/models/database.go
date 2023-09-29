package models

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type dbConfig struct {
	Host       string
	Port       string
	DBUser     string
	DBName     string
	DBPassword string
	SSLMode    string
}

func SetupDB() (*gorm.DB, error) {
	// Connect database
	db, err := ConnectDB()
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	log.Println("Successfully connected to database")

	// Migrate database
	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Successfully migrated database")

	return db, nil
}

func ConnectDB() (*gorm.DB, error) {
	dbConfig := newDBConfig()
	dsn := dbConfig.createDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}
	return db, nil
}

func newDBConfig() *dbConfig {
	return &dbConfig{
		Host:       os.Getenv("DB_HOST"),
		Port:       os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBName:     os.Getenv("DB_NAME"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		SSLMode:    os.Getenv("DB_SSL_MODE"),
	}
}

func (dbConfig *dbConfig) createDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.DBUser, dbConfig.DBName, dbConfig.DBPassword, dbConfig.SSLMode,
	)
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}
