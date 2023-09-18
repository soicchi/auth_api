package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDSN(t *testing.T) {
	dbConfig := &dbConfig{
		Host:     "host",
		Port:     "port",
		DBUser: "user",
		DBName: "database",
		DBPassword: "password",
		SSLMode: "disable",
	}
	dsn := dbConfig.createDSN()
	assert.Equal(t, "host=host port=port user=user dbname=database password=password sslmode=disable", dsn)
}

func TestValidateDBConfig(t *testing.T) {
	tests := []struct {
		name string
		dbConfig *dbConfig
		ErrMsg string
		wantErr bool
	}{
		{
			name: "valid",
			dbConfig: &dbConfig{
				Host: "host",
				Port: "port",
				DBUser: "user",
				DBName: "database",
				DBPassword: "password",
				SSLMode: "disable",
			},
			ErrMsg: "",
			wantErr: false,
		},
		{
			name: "invalid host",
			dbConfig: &dbConfig{
				Host: "",
				Port: "port",
				DBUser: "user",
				DBName: "database",
				DBPassword: "password",
				SSLMode: "disable",
			},
			ErrMsg: "DB_HOST is not set",
			wantErr: true,
		},
		{
			name: "invalid port",
			dbConfig: &dbConfig{
				Host: "host",
				Port: "",
				DBUser: "user",
				DBName: "database",
				DBPassword: "password",
				SSLMode: "disable",
			},
			ErrMsg: "DB_PORT is not set",
			wantErr: true,
		},
		{
			name: "invalid user",
			dbConfig: &dbConfig{
				Host: "host",
				Port: "port",
				DBUser: "",
				DBName: "database",
				DBPassword: "password",
				SSLMode: "disable",
			},
			ErrMsg: "DB_USER is not set",
			wantErr: true,
		},
		{
			name: "invalid dbname",
			dbConfig: &dbConfig{
				Host: "host",
				Port: "port",
				DBUser: "user",
				DBName: "",
				DBPassword: "password",
				SSLMode: "disable",
			},
			ErrMsg: "DB_NAME is not set",
			wantErr: true,
		},
		{
			name: "invalid password",
			dbConfig: &dbConfig{
				Host: "host",
				Port: "port",
				DBUser: "user",
				DBName: "database",
				DBPassword: "",
				SSLMode: "disable",
			},
			ErrMsg: "DB_PASSWORD is not set",
			wantErr: true,
		},
		{
			name: "invalid sslmode",
			dbConfig: &dbConfig{
				Host: "host",
				Port: "port",
				DBUser: "user",
				DBName: "database",
				DBPassword: "password",
				SSLMode: "invalid",
			},
			ErrMsg: "DB_SSL_MODE is invalid value",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.dbConfig.validateDBConfig()
			if tt.wantErr && err != nil {
				assert.Equal(t, tt.ErrMsg, err.Error())
			}
		})
	}
}
