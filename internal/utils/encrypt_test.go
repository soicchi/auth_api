package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptPassword(t *testing.T) {
	password := "password"
	hashedPassword, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEqual(t, password, hashedPassword)
}

func TestValidatePassword(t *testing.T) {
	password := "password"
	hashedPassword, err := HashPassword(password)
	assert.NoError(t, err)
	assert.True(t, ValidatePassword(hashedPassword, password))
}
