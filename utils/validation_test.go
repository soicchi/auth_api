package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/go-playground/validator/v10"
)

func TestNewCustomValidator(t *testing.T) {
	cv := NewCustomValidator()
	assert.NotNil(t, cv)
	assert.NotNil(t, cv.Validator)
	assert.IsType(t, &CustomValidator{}, cv)
	assert.IsType(t, &validator.Validate{}, cv.Validator)
}
