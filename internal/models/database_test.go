package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDSN(t *testing.T) {
	dbConfig := &dbConfig{
		Host:       "host",
		Port:       "port",
		DBUser:     "user",
		DBName:     "database",
		DBPassword: "password",
		SSLMode:    "disable",
	}
	dsn := dbConfig.createDSN()
	assert.Equal(t, "host=host port=port user=user dbname=database password=password sslmode=disable", dsn)
}
