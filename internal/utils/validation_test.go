package utils

import (
	"os"
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

func TestValidateENVVars(t *testing.T) {
	tests := []struct {
		name    string
		vars    []string
		wantErr bool
	}{
		{
			name:    "valid",
			vars:    []string{"TEST"},
			wantErr: false,
		},
		{
			name:    "invalid",
			vars:    []string{"TEST", "TEST2"},
			wantErr: true,
		},
		{
			name:    "empty",
			vars:    []string{},
			wantErr: false,
		},
		{
			name:    "nil",
			vars:    nil,
			wantErr: false,
		},
	}

	os.Setenv("TEST", "test")

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateENVVars(test.vars)
			if test.wantErr && err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
