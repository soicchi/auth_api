package utils

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestNewCustomValidator(t *testing.T) {
	cv := NewCustomValidator()
	assert.NotNil(t, cv)
	assert.NotNil(t, cv.Validator)
	assert.IsType(t, &CustomValidator{}, cv)
	assert.IsType(t, &validator.Validate{}, cv.Validator)
}
