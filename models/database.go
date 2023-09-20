package models

import (
	"fmt"
	"os"

	"golang.org/x/exp/slices"
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

func ConnectDB() (*gorm.DB, error) {
	dbConfig := newDBConfig()
	if err := dbConfig.validateDBConfig(); err != nil {
		return nil, fmt.Errorf("failed to validate DB config: %w", err)
	}
	
	dsn := dbConfig.createDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}
	return db, nil
}

func (dbConfig *dbConfig) validateDBConfig() error {
	sslModes := []string{"disable", "require"}

	if dbConfig.Host == "" {
		return fmt.Errorf("DB_HOST is not set")
	}
	if dbConfig.Port == "" {
		return fmt.Errorf("DB_PORT is not set")
	}
	if dbConfig.DBUser == "" {
		return fmt.Errorf("DB_USER is not set")
	}
	if dbConfig.DBName == "" {
		return fmt.Errorf("DB_NAME is not set")
	}
	if dbConfig.DBPassword == "" {
		return fmt.Errorf("DB_PASSWORD is not set")
	}
	if !slices.Contains(sslModes, dbConfig.SSLMode) {
		return fmt.Errorf("DB_SSL_MODE is invalid value")
	}
	return nil
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

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}