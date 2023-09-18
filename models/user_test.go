package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user := NewUser("email", "password")
	assert.Equal(t, "email", user.Email)
	assert.Equal(t, "password", user.Password)
}
